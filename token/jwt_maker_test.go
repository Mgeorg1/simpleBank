package token

import (
	"github.com/Mgeorg1/simpleBank/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, payload.Username, username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
	require.NotZero(t, payload.ID)
}

func TestJWTMakerExpired(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, payload, err := maker.CreateToken(username, util.DepositorRole, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
