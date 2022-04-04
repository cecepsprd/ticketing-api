package repository

import (
	"context"
	"database/sql"

	"github.com/cecepsprd/ticketing-api/model"
)

var (
	insertTransaction       = `INSERT INTO transaction (product_id, user_id, amount, status) VALUES (?,?,?,?)`
	updateTransaction       = `UPDATE transaction set payment_url=? WHERE id=?`
	updateTransactionStatus = `UPDATE transaction set status=? WHERE id=?`
	readTransactionByID     = `SELECT product_id, user_id, amount, status, payment_url, created_at, updated_at FROM product WHERE id=?`
)

type TransactionRepository interface {
	Create(context.Context, model.Transaction) (*model.Transaction, error)
	Update(context.Context, model.Transaction) error
	UpdateStatus(ctx context.Context, transactionID int64, status string) error
	ReadByID(ctx context.Context, transactionID int64) (*model.Transaction, error)
}

type mysqlTrxRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &mysqlTrxRepository{
		db: db,
	}
}

func (m *mysqlTrxRepository) Create(ctx context.Context, request model.Transaction) (*model.Transaction, error) {
	stmt, err := m.db.PrepareContext(ctx, insertTransaction)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(
		ctx,
		request.ProductID,
		request.UserID,
		request.Amount,
		request.Status,
	)
	if err != nil {
		return nil, err
	}

	request.ID, _ = result.LastInsertId()

	return &request, nil
}

func (repo *mysqlTrxRepository) Update(ctx context.Context, request model.Transaction) error {
	stmt, err := repo.db.PrepareContext(ctx, updateTransaction)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, request.PaymentURL, request.ID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *mysqlTrxRepository) UpdateStatus(ctx context.Context, transactionID int64, status string) error {
	stmt, err := repo.db.PrepareContext(ctx, updateTransactionStatus)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, status, transactionID)

	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlTrxRepository) ReadByID(ctx context.Context, transactionID int64) (*model.Transaction, error) {
	var transaction model.Transaction
	err := m.db.QueryRowContext(ctx, readProductByID, transactionID).Scan(
		&transaction.ID,
		&transaction.ProductID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.PaymentURL,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil
	}

	return &transaction, nil
}
