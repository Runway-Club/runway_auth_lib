package providers

import (
	"context"
	"github.com/Runway-Club/auth_lib/domain"
)

type DummyProvider struct {
	jwt domain.JwtGenerator
}

func NewDummyProvider(jwt domain.JwtGenerator) *DummyProvider {
	return &DummyProvider{jwt: jwt}
}

func (d DummyProvider) VerifyToken(ctx context.Context, token string) (uid string, claims map[string]interface{}, err error) {
	auth, claims, err := d.jwt.VerifyToken(token)
	if err != nil {
		return "", nil, err
	}
	return auth.Id, claims, nil
}
