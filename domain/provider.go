package domain

import "context"

type Provider interface {
	VerifyToken(ctx context.Context, token string) (uid string, claims map[string]interface{}, err error)
	Delete(ctx context.Context, uid string) error
}
