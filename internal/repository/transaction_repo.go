package repository

import (
	"AccountManagementSystem/internal/models"
	"context"
	"database/sql"
)

type TransactionRepo struct{ DB *sql.DB }

func NewTransactionRepo(db *sql.DB) *TransactionRepo { return &TransactionRepo{DB: db} }

func (r *TransactionRepo) Insert(ctx context.Context, tx *sql.Tx, t models.Transaction) (models.Transaction, error) {
	row := tx.QueryRowContext(ctx, `INSERT INTO transactions (account_id, type, amount) VALUES ($1,$2,$3) RETURNING id, account_id, type, amount, created_at`, t.AccountID, t.Type, t.Amount)
	var out models.Transaction
	if err := row.Scan(&out.ID, &out.AccountID, &out.Type, &out.Amount, &out.CreatedAt); err != nil {
		return out, err
	}
	return out, nil
}

func (r *TransactionRepo) ListByAccount(ctx context.Context, accountID int64, limit int) ([]models.Transaction, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, account_id, type, amount, created_at FROM transactions WHERE account_id=$1 ORDER BY created_at DESC LIMIT $2`, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]models.Transaction, 0)
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.Type, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}
