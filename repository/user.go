package repository

import (
	"context"

	"database/sql"

	"github.com/cecepsprd/ticketing-api/model"
)

var (
	readUserByID       = `SELECT id, username, email, password, phone, address, created_at, updated_at FROM user WHERE id = ?`
	insertUser         = `INSERT INTO user (username, password, email, phone, address, roles) VALUES (?,?,?,?,?,?)`
	readAllUser        = `SELECT id, username, email, password, phone, address, roles FROM user`
	updateUser         = `UPDATE user set username=?, email=?, password=?, phone=?, address=?, roles=? WHERE id = ?`
	deleteUser         = `DELETE FROM user WHERE id = ?`
	readUserByUsername = `SELECT id, username, email, password, phone, address, roles, created_at, updated_at FROM user WHERE username = ?`
	checkDuplicateUser = `SELECT count(1) FROM user WHERE username=? or email=? or phone=?`
)

type UserRepository interface {
	Create(context.Context, model.User) error
	Read(context.Context) ([]model.User, error)
	Update(ctx context.Context, request model.User) error
	Delete(ctx context.Context, userid int64) error
	ReadByID(ctx context.Context, userid int64) (user model.User, err error)
	ReadByUsername(ctx context.Context, username string) (user model.User, err error)
	CountUser(ctx context.Context, request model.User) (int32, error)
}

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{
		db: db,
	}
}

func (m *mysqlUserRepository) ReadByID(ctx context.Context, userid int64) (model.User, error) {
	var user model.User
	err := m.db.QueryRowContext(ctx, readUserByID, userid).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Address,
		&user.Roles,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return model.User{}, nil
	}

	return user, nil
}

func (m *mysqlUserRepository) Create(ctx context.Context, request model.User) error {
	stmt, err := m.db.PrepareContext(ctx, insertUser)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, request.Username, request.Password, request.Email, request.Phone, request.Address, request.Roles)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mysqlUserRepository) Read(ctx context.Context) (response []model.User, err error) {
	rows, err := repo.db.QueryContext(ctx, readAllUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Phone,
			&user.Address,
			&user.Roles,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, user)
	}

	return response, nil
}

func (repo *mysqlUserRepository) Update(ctx context.Context, request model.User) error {
	stmt, err := repo.db.PrepareContext(ctx, updateUser)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, request.Username, request.Password, request.Email, request.Phone, request.Address, request.Roles, request.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mysqlUserRepository) Delete(ctx context.Context, userid int64) error {
	stmt, err := repo.db.PrepareContext(ctx, deleteUser)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userid)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlUserRepository) ReadByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := m.db.QueryRowContext(ctx, readUserByUsername, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Address,
		&user.Roles,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return model.User{}, nil
	}

	return user, nil
}

func (m *mysqlUserRepository) CountUser(ctx context.Context, request model.User) (total int32, err error) {
	err = m.db.QueryRowContext(ctx, checkDuplicateUser, request.Username, request.Email, request.Phone).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
