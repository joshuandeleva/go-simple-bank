package db
import (
	"context"
	"testing"
	"time"

	"github.com/joshuandeleva/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createARandomAccount(t)
	
}

func createARandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword: "secret",
		FullName: util.RandomOwner(),
		Email: util.RandomEamil(),
		
	}
	user, err := testQuerries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user) // checks if account is not empty

	require.Equal(t, arg.Username, user.Username) // checks if owner is equal to arg owner
	require.Equal(t, arg.FullName, user.FullName) // checks if balance is equal to arg balance
	require.Equal(t, arg.Email, user.Email) // checks if currency is equal to arg currency
	require.Equal(t, user.HashedPassword, user.HashedPassword) // checks if password is equal to secret

	require.NotZero(t, user.CreatedAt) // checks if account created at is not zero
	require.True(t, user.PasswordChangedAt.IsZero()) // checks if account id is not zero

	return user
} 

func TestGetUser(t *testing.T)  {
	// create account
	user1 := createARandomUser(t) // creates a random account
	// get account
	foundUser, err := testQuerries.GetUser(context.Background(),user1.Username)

	require.NoError(t, err) // checks if there is no error
	require.NotEmpty(t, foundUser) // checks if account is not empty

	require.Equal(t, user1.Username, foundUser.Username) // checks if owner is equal to account1 owner
	require.Equal(t, user1.FullName, foundUser.FullName) // checks if balance is equal to account1 balance
	require.Equal(t, user1.Email, foundUser.Email) // checks if currency is equal to account1 currency
	require.Equal(t, user1.HashedPassword, foundUser.HashedPassword) // checks if password is equal to account1 password
	require.WithinDuration(t, user1.CreatedAt, foundUser.CreatedAt, time.Second) // checks if account1 created at is equal to found account created at

}
