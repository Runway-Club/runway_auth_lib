package domain

import "errors"

type JwtGenerator interface {
	GenerateToken(auth *Auth, payload map[string]interface{}) (string, error)
	VerifyToken(token string) (*Auth, map[string]interface{}, error)
}

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidIssuer = errors.New("invalid issuer")
	ErrExpiredToken  = errors.New("expired token")
)
