package providers

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type GoogleProvider struct {
	client *auth.Client
}

func (g *GoogleProvider) VerifyToken(ctx context.Context, token string) (uid string, claims map[string]interface{}, err error) {
	// verify token
	verifiedToken, err := g.client.VerifyIDToken(ctx, token)
	if err != nil {
		return "", nil, err
	}
	return verifiedToken.UID, verifiedToken.Claims, nil
}

func NewGoogleProvider(ctx context.Context, configFileName string) *GoogleProvider {
	// initialize firebase app
	opt := option.WithCredentialsFile(configFileName)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		panic(err)
	}
	// initialize firebase auth client
	client, err := app.Auth(ctx)
	if err != nil {
		panic(err)
	}
	return &GoogleProvider{
		client: client,
	}
}
