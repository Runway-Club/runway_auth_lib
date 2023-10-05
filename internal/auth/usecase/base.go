package usecase

import (
	"context"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
)

type AuthUseCase struct {
	repo           domain.AuthRepository
	passwordPolicy string
	hashCost       string
	jwt            domain.JwtGenerator
	defaultRoleId  string
}

func (a *AuthUseCase) SignInWithProvider(ctx context.Context, provider domain.Provider, token string) (genToken *domain.Token, err error) {
	uid, claims, err := provider.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}
	// get by username
	user, err := a.repo.GetById(ctx, uid)
	if err != nil {
		return nil, domain.ErrAuthNotFound
	}
	// generate token
	generatedToken, err := a.jwt.GenerateToken(user, claims)
	if err != nil {
		return nil, err
	}
	if user.RoleId == "" {
		user.RoleId = a.defaultRoleId
		err = a.repo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return &domain.Token{
		Jwt:    generatedToken,
		Id:     user.Id,
		UserId: user.Id,
		RoleId: user.RoleId,
	}, nil
}

func NewAuthUseCase(repo domain.AuthRepository, jwt domain.JwtGenerator) *AuthUseCase {
	usecase := &AuthUseCase{
		repo:           repo,
		passwordPolicy: viper.GetString("runway_auth.password.policy"),
		hashCost:       viper.GetString("runway_auth.password.cost"),
		defaultRoleId:  viper.GetString("runway_auth.default_role_id"),
		jwt:            jwt,
	}
	// init static users, omit error because it's okay if it's already exist
	for _, user := range repo.GetStaticUserMap(context.Background()) {
		err := usecase.SignUp(context.Background(), user)
		if err != nil {
			log.Print(err)
		}
	}
	return usecase
}

func (a *AuthUseCase) SignUp(ctx context.Context, auth *domain.Auth) error {
	if a.passwordPolicy == "" {
		a.passwordPolicy = "level1"
	}
	// check password policy
	if a.passwordPolicy == "level1" {
		// minimum 8 characters
		if len(auth.Password) < 8 {
			return domain.ErrInvalidPassword
		}
	}
	if a.passwordPolicy == "level2" {
		//  minimum 8 any characters and contain at least one number
		regex := regexp.MustCompile(`[0-9]`)
		if len(auth.Password) < 8 || !regex.MatchString(auth.Password) {
			return domain.ErrInvalidPassword
		}
	}
	if a.passwordPolicy == "level3" {
		// minimum 8 any characters and contain at least one number and one uppercase letter and one special character
		regex := regexp.MustCompile(`[0-9]`)
		regex2 := regexp.MustCompile(`[A-Z]`)
		regex3 := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)

		if len(auth.Password) < 8 || !regex.MatchString(auth.Password) || !regex2.MatchString(auth.Password) || !regex3.MatchString(auth.Password) {
			return domain.ErrInvalidPassword
		}
	}
	// look for existing username
	found, err := a.repo.GetByUsername(ctx, auth.Username)
	if err == nil || found != nil {
		return domain.ErrUsernameExist
	}
	// hash password
	hashCost := bcrypt.DefaultCost
	if a.hashCost == "min" {
		hashCost = bcrypt.MinCost
	}
	if a.hashCost == "max" {
		hashCost = bcrypt.MaxCost
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(auth.Password), hashCost)
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

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (token *domain.Token, err error) {

	// get by username
	user, err := a.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrAuthNotFound
	}
	if user.RoleId == "" {
		user.RoleId = a.defaultRoleId
		err = a.repo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hpassword), []byte(password))
	if err != nil {
		return nil, domain.ErrPasswordNotMatch
	}
	// generate token
	generatedToken, err := a.jwt.GenerateToken(user, map[string]interface{}{
		"username": user.Username,
		"id":       user.Id,
		"role_id":  user.RoleId,
	})
	if err != nil {
		return nil, err
	}
	return &domain.Token{
		Jwt:    generatedToken,
		Id:     user.Id,
		UserId: user.Id,
		RoleId: user.RoleId,
	}, nil
}
