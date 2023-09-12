package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TranferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

//NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

//execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	err = fn(New(tx))

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

//TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//TransferTx performs a money transfer from one account to the other.
//It creates a transfer record, add account entries, add update accounts' balance within a single database transaction
func (store *SQLStore) TranferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = addEntry(q, ctx, arg.FromAccountID, -arg.Amount)

		if err != nil {
			return err
		}

		result.ToEntry, err = addEntry(q, ctx, arg.ToAccountID, +arg.Amount)

		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {

			result.FromAccount, result.ToAccount, err = addMoney(q, ctx, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

			if err != nil {
				return err
			}

		} else {

			result.ToAccount, result.FromAccount, err = addMoney(q, ctx, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err
}

func addMoney(
	q *Queries,
	ctx context.Context,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {

	account1, err = updateBalance(q, ctx, accountID1, amount1)

	if err != nil {
		return
	}

	account2, err = updateBalance(q, ctx, accountID2, amount2)

	return account1, account2, err

}

func addEntry(q *Queries, ctx context.Context, accountID int64, amount int64) (Entry, error) {

	return q.CreateEntry(ctx, CreateEntryParams{
		AccountID: accountID,
		Amount:    amount,
	})
}

func updateBalance(q *Queries, ctx context.Context, accountID int64, amount int64) (Account, error) {
	return q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID,
		Amount: amount,
	})
}
