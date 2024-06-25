package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)

}

// SQL store provides all functions to execute querries and transactions
type SQLStore struct {
	*Queries // embeds *Queries
	db       *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction

func (store *SQLStore) execTx(ctx context.Context, fun func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // nil means default options

	if err != nil {
		return err
	}

	q := New(tx) // creates a new Queries object with the transaction
	err = fun(q) // executes the function
	if err != nil {

		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("txn error: %w , rb error: %w", err, rbError)
		}
		return err // returns the original error

	}
	return tx.Commit() // commits the transaction

}

//  perfomr money transfer from one account to another

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entery   `json:"from_entry"`
	ToEntry     Entery   `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// creates a transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))

		if err != nil {
			return err
		}

		// creates an entry for the from account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// creates an entry for the to account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// to avoid locks use conditional

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount , err =  addMoney(ctx,q,arg.FromAccountID,-arg.Amount,arg.ToAccountID,arg.Amount)

		} else {
			result.ToAccount, result.FromAccount , err =  addMoney(ctx,q,arg.ToAccountID,arg.Amount,arg.FromAccountID,-arg.Amount)

		}
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}


func addMoney(ctx context.Context,q *Queries,accountID1 int64,amount1 int64,accountID2 int64,amount2 int64)(account1 Account, account2 Account, err error) {
	account1 , err = q.UpdateAccountBalance(ctx ,UpdateAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return 
	}
	account2 , err = q.UpdateAccountBalance(ctx ,UpdateAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return 
	}
	return 
}