package db

import (
	"context"
	"testing"

	"github.com/joshuandeleva/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entery {

	account1 := createARandomAccount(t) // creates a random account
	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQuerries.CreateEntry(context.Background(), arg)
	require.NoError(t, err) // checks if there is an error
	require.NotEmpty(t, entry) // checks if entry is not empty

	require.Equal(t, arg.AccountID, entry.AccountID) // checks if account id is equal to arg account id
	require.Equal(t, arg.Amount, entry.Amount) // checks if amount is equal to arg amount
	require.NotZero(t, entry.ID) // checks if entry id is not zero
	require.NotZero(t, entry.CreatedAt) // checks if entry created at is not zero
	require.Equal(t, account1.ID, entry.AccountID) // checks if account id is equal to account1 id
	
	return entry
}

func TestCreatingEntry(t *testing.T){
	 createRandomEntry(t) // creates a random entry
}