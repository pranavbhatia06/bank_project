package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbERR := tx.Rollback(); rbERR != nil {
			return fmt.Errorf("tx Err %v ,rb Err: %v", err, rbERR)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	FromAccountId Account  `json:"from_account_id"`
	ToAccountId   Account  `json:"to_account_id"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
	Transfer      Transfer `json:"transfer"`
}

func (s *Store) TransferTx(ctx context.Context, transferParam TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	debugLog := ctx.Value(struct{}{})
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		fmt.Println(debugLog, "Create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: transferParam.FromAccountID,
			ToAccountID:   transferParam.ToAccountID,
			Amount:        transferParam.Amount,
		})

		if err != nil {
			return err
		}
		fmt.Println(debugLog, "Create Entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transferParam.ToAccountID,
			Amount:    transferParam.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(debugLog, "Create Entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transferParam.FromAccountID,
			Amount:    -(transferParam.Amount),
		})
		if err != nil {
			return err
		}
		if transferParam.FromAccountID < transferParam.ToAccountID {
			err = s.updateAccountsBalance(ctx, transferParam.FromAccountID, transferParam.ToAccountID, -(transferParam.Amount), transferParam.Amount, &result)
		} else {
			err = s.updateAccountsBalance(ctx, transferParam.ToAccountID, transferParam.FromAccountID, transferParam.Amount, -(transferParam.Amount), &result)
		}
		if err != nil {
			return err
		}

		return nil

	})

	return result, err
}
func (s *Store) updateAccountsBalance(ctx context.Context, account1Id int64, account2Id int64, amount1 int64, amount2 int64, result *TransferTxResult) error {
	debugLog := ctx.Value(struct{}{})
	fmt.Println(debugLog, "Add Amount Account 1")
	account1, err := s.AddAmountToAccount(ctx, AddAmountToAccountParams{
		Amount: amount1,
		ID:     account1Id,
	})
	if err != nil {
		return err
	}
	fmt.Println(debugLog, "Add Amount Account 2")
	account2, err := s.AddAmountToAccount(ctx, AddAmountToAccountParams{
		Amount: amount2,
		ID:     account2Id,
	})
	if err != nil {
		return err
	}
	if amount1 < amount2 {
		result.FromAccountId = account1
		result.ToAccountId = account2
	} else {
		result.ToAccountId = account1
		result.FromAccountId = account2
	}
	return nil

}
