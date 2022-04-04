package repository

import (
	"context"
	"database/sql"

	"github.com/cecepsprd/ticketing-api/model"
)

var (
	insertProduct   = `INSERT INTO product (name, description, price, stock, image_url, start_date, end_date) VALUES (?,?,?,?,?,?,?)`
	updateProduct   = `UPDATE product SET (name, description, price, stock, image_url, start_date, end_date) VALUES (?,?,?,?,?,?,?)`
	readAllProduct  = `SELECT id, name, description, price, stock, image_url, start_date, end_date, created_at, updated_at FROM product`
	deleteProduct   = `DELETE FROM product WHERE id = ?`
	readProductByID = `SELECT id, name, description, price, stock, image_url, start_date, end_date, created_at, updated_at FROM product WHERE id=?`
	updateStock     = `UPDATE product SET stock=(stock+?) WHERE id=?`
)

type ProductRepository interface {
	Create(ctx context.Context, product model.Product) error
	Read(context.Context) ([]model.Product, error)
	Update(ctx context.Context, product model.Product) error
	Delete(ctx context.Context, productID int64) error
	ReadByID(ctx context.Context, id int64) (*model.Product, error)
	UpdateStock(ctx context.Context, productID int64, newStock int64) error
}

type mysqlProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &mysqlProductRepository{
		db: db,
	}
}

func (repo *mysqlProductRepository) Create(ctx context.Context, request model.Product) error {
	stmt, err := repo.db.PrepareContext(ctx, insertProduct)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		request.Name,
		request.Description,
		request.Price,
		request.Stock,
		request.ImageURL,
		request.StartDate,
		request.EndDate,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mysqlProductRepository) Read(ctx context.Context) (response []model.Product, err error) {
	rows, err := repo.db.QueryContext(ctx, readAllProduct)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.ImageURL,
			&p.StartDate,
			&p.EndDate,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, p)
	}

	return response, nil
}

func (repo *mysqlProductRepository) Update(ctx context.Context, request model.Product) error {
	stmt, err := repo.db.PrepareContext(ctx, updateProduct)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		request.Name,
		request.Description,
		request.Price,
		request.Stock,
		request.ImageURL,
		request.StartDate,
		request.EndDate,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mysqlProductRepository) Delete(ctx context.Context, userid int64) error {
	stmt, err := repo.db.PrepareContext(ctx, deleteProduct)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userid)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlProductRepository) ReadByID(ctx context.Context, productID int64) (*model.Product, error) {
	var p model.Product
	err := m.db.QueryRowContext(ctx, readProductByID, productID).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.ImageURL,
		&p.StartDate,
		&p.EndDate,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil
	}

	return &p, nil
}

func (repo *mysqlProductRepository) UpdateStock(ctx context.Context, productID int64, newStock int64) error {
	stmt, err := repo.db.PrepareContext(ctx, updateStock)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, newStock, productID)
	if err != nil {
		return err
	}

	return nil
}
