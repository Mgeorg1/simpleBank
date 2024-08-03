package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTMaker struct {
	secretKey string
}

const MinSecretKeySize = 32

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d", MinSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		var validErr *jwt.ValidationError
		ok := errors.As(err, &validErr)
		if ok && errors.Is(validErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
