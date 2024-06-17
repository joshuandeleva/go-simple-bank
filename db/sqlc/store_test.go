package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createARandomAccount(t)
	account2 := createARandomAccount(t)


	// run a concurrent transfer transaction

	n := 5

	amount := int64(n)

	// use channels to send and receive the result of the transfer transaction
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background() // add transaction name to context
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			// send the result and error to the channel
			errs <- err
			results <- result
		}()
	}

	// wait for all the goroutines to finish
	existed := make(map[int]bool) // map to check if the transfer has been created

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check parameters in the result data from chanel
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check if the record has been created in the database

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check account entries of result
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount) // check if the amount is negative money is going out of the account
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt) // check if the entry has been created

		// get the account entry from the database
		_, err = store.GetEntry(context.Background(), fromEntry.ID) // check if the entry has been created
		require.NoError(t, err)                                     // check if there is no error

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount) // check if the amount is positive money is going into the account
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt) // check if the entry has been created

		// get the account entry from the database
		_, err = store.GetEntry(context.Background(), toEntry.ID) // check if the entry has been created
		require.NoError(t, err)                                   // check if there is no error

		// check if the account balances have been updated

		// check account

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check if the account balances have been updated
		// logs


		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2) // check if the account balances have been updated
		require.True(t, diff1 > 0)  // shd be positive
		require.True(t, diff1 % amount == 0)   // 1 * amount , 2 * amount, 3 * amount, 4 * amount, 5 * amount, ... n * amount

		k := int(diff1 / amount) // check if the amount is divisible by the difference
		require.True(t, k >= 1 && k <= n) // check if the amount is divisible by the difference

		require.NotContains(t, existed, k) // check if the transfer has been created

		existed[k] = true // add the transfer to the existed map

	}

	// check update final balance

	updateAccount1 , err := testQuerries.GetAccount(context.Background(), account1.ID) // get the account from the database
	require.NoError(t, err)

	updateAccount2 , err := testQuerries.GetAccount(context.Background(), account2.ID) // get the account from the database
	require.NoError(t, err)


	require.Equal(t, account1.Balance -  int64(n) * amount , updateAccount1.Balance) // check if the account balances have been updated
	require.Equal(t, account2.Balance +  int64(n) * amount , updateAccount2.Balance) // check if the account balances have been updated
 

}
func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createARandomAccount(t)
	account2 := createARandomAccount(t)


	// run a concurrent transfer transaction

	n := 10

	amount := int64(10)

	// use channels to send and receive the result of the transfer transaction
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i % 2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			ctx := context.Background() // add transaction name to context
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			// send the result and error to the channel
			errs <- err
		}()
	}

	// wait for all the goroutines to finish

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	// check update final balance

	updateAccount1 , err := testQuerries.GetAccount(context.Background(), account1.ID) // get the account from the database
	require.NoError(t, err)

	updateAccount2 , err := testQuerries.GetAccount(context.Background(), account2.ID) // get the account from the database
	require.NoError(t, err)


	require.Equal(t, account1.Balance , updateAccount1.Balance) // check if the account balances have been updated
	require.Equal(t, account2.Balance , updateAccount2.Balance) // check if the account balances have been updated
 

}
