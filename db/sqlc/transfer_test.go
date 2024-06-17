package db

// import (
// 	"context"
// 	"testing"

// 	"github.com/joshuandeleva/simplebank/util"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomTransfer(t *testing.T) Transfer {
// 	account1 := createARandomAccount(t) // creates a random account
// 	account2 := createARandomAccount(t) // creates a random account

// 	// the amount to transfer should be less than the balance of the account and greater than zero
// 	n := util.RandomInt(0, account1.Balance) // generates a random number between 0 and the balance of the account
// 	if n == 0 { // if the random number is 0
// 		n = 1 // then set the random number to 1
// 	}

// 	arg := CreateTransferParams{
// 		FromAccountID: account1.ID,
// 		ToAccountID: account2.ID,
// 		Amount: n,
// 	}
// 	transfer, err := testQuerries.CreateTransfer(context.Background(), arg) // creates a transfer
// 	require.NoError(t, err) // checks if there is an error
// 	require.NotEmpty(t, transfer) // checks if transfer is not empty
// 	require.Equal(t, arg.FromAccountID, transfer.FromAccountID) // checks if the from account id is equal to the transfer from account id
// 	require.Equal(t, arg.ToAccountID, transfer.ToAccountID) // checks if the to account id is equal to the transfer to account id
// 	require.Equal(t, arg.Amount, transfer.Amount) // checks if the amount is equal to the transfer Amount
// 	require.NotZero(t, transfer.ID) // checks if the transfer id is not zero
// 	require.NotZero(t, transfer.CreatedAt) // checks if the transfer created at is not zero
// 	require.Equal(t, account1.ID, transfer.FromAccountID) // checks if the from account id is equal to the account1 id
// 	require.Equal(t, account2.ID, transfer.ToAccountID) // checks if the to account id is equal to the account2 id
// 	require.NotZero(t, transfer.CreatedAt) // checks if the transfer created at is not zero
// 	return transfer // returns the transfer
// }

// func TestCreateTransfer(t *testing.T) {
// 	createRandomTransfer(t) // creates a random transfer
// }