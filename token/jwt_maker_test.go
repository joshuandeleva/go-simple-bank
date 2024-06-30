package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joshuandeleva/simplebank/util"
	"github.com/stretchr/testify/require"
)


func TestJwtMaker(t *testing.T){
	maker , err := NewJWTMaker(util.RandomString(32))
	require.NoError(t , err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token , err := maker.CreateToken(username , duration)
	fmt.Println("Generated Token:", token)
	require.NoError(t , err)
	require.NotEmpty(t , token)

	payload , err := maker.VerifyToken(token)
	require.NoError(t , err)
	require.NotEmpty(t , payload)

	require.NotZero(t , payload.ID)
	require.Equal(t , username , payload.Username)
	require.WithinDuration(t , issuedAt , payload.IssuedAt , time.Second) // time.Second is the tolerance
	require.WithinDuration(t , expiredAt , payload.ExpiredAt , time.Second)
}

func TestExpiredJwtToken(t *testing.T){
	maker , err := NewJWTMaker(util.RandomString(32))
	require.NoError(t , err)

	token , err := maker.CreateToken(util.RandomOwner() , -time.Minute) // negative duration
	require.NoError(t , err)
	require.NotEmpty(t , token)

	payload , err := maker.VerifyToken(token)
	require.Error(t , err)
	require.EqualError(t , err , ErrExpiredToken.Error())
	require.Nil(t , payload)

}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err) // no error

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err) // no error

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err) // error
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)	
}