package user

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	const op = "user.Repository.CreateUser"
	var lastID int

	query := "INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&lastID)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	user.ID = int64(lastID)

	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const op = "user.Repository.GetUserByEmail"
	u := User{}

	query := "SELECT id, email, encrypted_password, created_at FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	return &u, nil
}
