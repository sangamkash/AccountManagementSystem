package repository

import (
	"AccountManagementSystem/internal/models"
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("not found")

type AccountRepo struct{ DB *sql.DB }

func NewAccountRepo(db *sql.DB) *AccountRepo { return &AccountRepo{DB: db} }

func (r *AccountRepo) Create(ctx context.Context, username string, initial float64) (*models.Account, error) {
	row := r.DB.QueryRowContext(ctx, `INSERT INTO accounts (username, balance) VALUES ($1,$2) RETURNING id, username, balance, created_at`, username, initial)
	var a models.Account
	if err := row.Scan(&a.ID, &a.Username, &a.Balance, &a.CreatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountRepo) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	var a models.Account
	row := r.DB.QueryRowContext(ctx, `SELECT id, username, balance, created_at FROM accounts WHERE id=$1`, id)
	if err := row.Scan(&a.ID, &a.Username, &a.Balance, &a.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

// GetForUpdate selects account row FOR UPDATE inside a transaction
func (r *AccountRepo) GetForUpdate(ctx context.Context, tx *sql.Tx, id int64) (models.Account, error) {
	var a models.Account
	row := tx.QueryRowContext(ctx, `SELECT id, username, balance, created_at FROM accounts WHERE id=$1 FOR UPDATE`, id)
	if err := row.Scan(&a.ID, &a.Username, &a.Balance, &a.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return a, ErrNotFound
		}
		return a, err
	}
	return a, nil
}

func (r *AccountRepo) UpdateBalance(ctx context.Context, tx *sql.Tx, id int64, newBalance float64) error {
	_, err := tx.ExecContext(ctx, `UPDATE accounts SET balance=$1 WHERE id=$2`, newBalance, id)
	return err
}
