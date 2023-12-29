// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: transfers.sql

package sqlc

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id, to_account_id ,amount
) VALUES (
             $1, $2, $3
         )
    RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listAccountCredits = `-- name: ListAccountCredits :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE to_account_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListAccountCredits(ctx context.Context, toAccountID int64) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listAccountCredits, toAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAccountTransactions = `-- name: ListAccountTransactions :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE from_account_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListAccountTransactions(ctx context.Context, fromAccountID int64) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listAccountTransactions, fromAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransactionFromOneAccountToOther = `-- name: ListTransactionFromOneAccountToOther :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE to_account_id = $1 AND from_account_id = $2
ORDER BY created_at DESC
`

type ListTransactionFromOneAccountToOtherParams struct {
	ToAccountID   int64
	FromAccountID int64
}

func (q *Queries) ListTransactionFromOneAccountToOther(ctx context.Context, arg ListTransactionFromOneAccountToOtherParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransactionFromOneAccountToOther, arg.ToAccountID, arg.FromAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
