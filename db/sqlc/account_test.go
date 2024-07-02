package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/joshuandeleva/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createARandomAccount(t)
	
}

func createARandomAccount(t *testing.T) Account {
	user := createARandomUser(t) // creates a random user
	arg := CreateAccountParams{
		Owner: user.Username ,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQuerries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account) // checks if account is not empty

	require.Equal(t, arg.Owner, account.Owner) // checks if owner is equal to arg owner
	require.Equal(t, arg.Balance, account.Balance) // checks if balance is equal to arg balance
	require.Equal(t, arg.Currency, account.Currency) // checks if currency is equal to arg currency

	require.NotZero(t, account.ID) // checks if account id is not zero
	require.NotZero(t, account.CreatedAt) // checks if account created at is not zero

	return account
} 

func TestGetAccount(t *testing.T)  {
	// create account
	account1 := createARandomAccount(t) // creates a random account
	// get account
	account2, err := testQuerries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err) // checks if there is no error
	require.NotEmpty(t, account2) // checks if account is not empty

	require.Equal(t, account1.Owner, account2.Owner) // checks if owner is equal to account1 owner
	require.Equal(t, account1.ID, account2.ID) // checks if balance is equal to account1 balance
	require.Equal(t, account1.Currency, account2.Currency) // checks if currency is equal to account1 currency
	require.Equal(t, account1.Balance, account2.Balance) // checks if balance is equal to account1 balance
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // checks if created at is within 0 duration


}

func TestUpdateAccount(t *testing.T) {
	// create account
	account1 := createARandomAccount(t) // creates a random account
	//fields to update
	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: account1.Balance,
	}
	// update account
	account2, err := testQuerries.UpdateAccount(context.Background(), arg) // updates account
	require.NoError(t, err) // checks if there is no error
	require.NotEmpty(t, account2) // checks if account is not empty

	require.Equal(t, account1.Owner, account2.Owner) // checks if owner is equal to account1 owner
	require.Equal(t, account1.ID, account2.ID) // checks if balance is equal to account1 balance
	require.Equal(t, account1.Currency, account2.Currency) // checks if currency is equal to account1 currency
	require.Equal(t, arg.Balance, account2.Balance) // checks if balance is equal to account1 balance
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // checks if created at is within 0 duration
	
}

func TestDeleteAccount(t *testing.T) {
	// create account
	account1 := createARandomAccount(t) // creates a random account

	// delete account
	err := testQuerries.DeleteAccount(context.Background(), account1.ID) // deletes account
	require.NoError(t, err) // checks if there is no error

	// get account
	account2, err := testQuerries.GetAccount(context.Background(), account1.ID) // gets account
	require.Error(t, err) // checks if there is an error
	require.EqualError(t, err, sql.ErrNoRows.Error()) // checks if error is sql no rows in result set
	require.Empty(t, account2) // checks if account is empty

}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	// use a for loop to create 10 accounts
	for i := 0; i < 10; i++ {
		lastAccount = createARandomAccount(t) // creates a random account
	}
	// list account args i.e limit and offset
	arg := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}
	// list accounts
	accounts, err := testQuerries.ListAccounts(context.Background() ,arg) // lists accounts
	require.NoError(t, err) // checks if there is no error
	require.NotEmpty(t, accounts) 

	// loop through accounts and check if each account is not empty
	for _, account := range accounts {
		require.NotEmpty(t, account) // checks if account is not empty
		require.Equal(t, arg.Owner, account.Owner) // checks if owner is equal to arg owner
	}
}