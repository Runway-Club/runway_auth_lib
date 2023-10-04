package usecase

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"runwayclub.dev/auth/domain"
	"runwayclub.dev/auth/internal/auth/repo"
)

type AuthUseCase struct {
	repo *repo.AuthRepository
}

func NewAuthUseCase(repo *repo.AuthRepository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

func (a *AuthUseCase) SignUp(ctx context.Context, auth *domain.Auth) error {
	// look for existing username
	found, err := a.repo.GetByUsername(ctx, auth.Username)
	if err != nil {
		return domain.ErrInternal
	}
	if found != nil {
		return domain.ErrUsernameExist
	}
	// hash password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.ErrInternal
	}
	auth.Hpassword = string(hashedPassword)

	// create new auth
	err = a.repo.Create(ctx, auth)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (*domain.Auth, error) {
	// get by username
	user, err := a.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrAuthNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hpassword), []byte(password))
	if err != nil {
		return nil, domain.ErrPasswordNotMatch
	}
}
