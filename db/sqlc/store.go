package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rollBackError := tx.Rollback()
		if rollBackError != nil {
			fmt.Errorf("transaction error: %v, rollback error: %v", err, rollBackError)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, params CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, params)
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, err = q.GetAccount(ctx, params.FromAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.GetAccount(ctx, params.ToAccountID)
		if err != nil {
			return err
		}

		fromResultBalance := result.FromAccount.Balance - params.Amount
		toResultBalance := result.ToAccount.Balance + params.Amount
		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      params.FromAccountID,
			Balance: fromResultBalance,
		})
		if err != nil {
			return err
		}

		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      params.ToAccountID,
			Balance: toResultBalance,
		})

		// TODO: take care about deadlock

		return err
	})

	return result, err
}
