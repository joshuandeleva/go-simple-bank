package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/joshuandeleva/simplebank/util"
	"github.com/stretchr/testify/require"
)


func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	fmt.Println("Generated Token:", token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second) // time.Second is the tolerance
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T){
	maker , err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t , err)

	token , err := maker.CreateToken(util.RandomOwner() , -time.Minute) // negative duration
	require.NoError(t , err)
	require.NotEmpty(t , token)

	payload , err := maker.VerifyToken(token)
	require.Error(t , err)
	require.EqualError(t , err , ErrExpiredToken.Error())
	require.Nil(t , payload)

}
