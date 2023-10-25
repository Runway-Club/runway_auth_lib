package usecase

import (
	"context"
	"fmt"
	"github.com/Runway-Club/auth_lib/common"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type AuthUseCase struct {
	repo           domain.AuthRepository
	passwordPolicy string
	hashCost       string
	jwt            domain.JwtGenerator
	defaultRoleId  string
}

func (a *AuthUseCase) ChangePassword(ctx context.Context, uid, oldPassword, newPassword string) error {
	// get by username
	user, err := a.repo.GetById(ctx, uid)
	if err != nil {
		return domain.ErrAuthNotFound
	}

	// check password
	errPasswordPolicy := common.CheckPasswordPolicy(newPassword, a.passwordPolicy)
	if errPasswordPolicy != nil {
		return errPasswordPolicy
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hpassword), []byte(oldPassword))
	if err != nil {
		return domain.ErrPasswordNotMatch
	}
	// hash password
	hashedPassword, err := common.GeneratePassword(newPassword, a.hashCost)
	if err != nil {
		return err
	}
	user.Hpassword = hashedPassword
	err = a.repo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) ChangeRole(ctx context.Context, uid, roleId string) error {
	// get by username
	user, err := a.repo.GetById(ctx, uid)
	if err != nil {
		return domain.ErrAuthNotFound
	}
	user.RoleId = roleId
	err = a.repo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) GetByUsername(ctx context.Context, username string) (*domain.Auth, error) {
	return a.repo.GetByUsername(ctx, username)
}

func (a *AuthUseCase) GetById(ctx context.Context, id string) (*domain.Auth, error) {
	return a.repo.GetById(ctx, id)
}

func (a *AuthUseCase) List(ctx context.Context, opt *common.QueryOpts) (*common.ListResult[*domain.Auth], error) {
	return a.repo.List(ctx, opt)
}

func (a *AuthUseCase) Verify(ctx context.Context, token string) (auth *domain.Auth, err error) {
	auth, _, err = a.jwt.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (a *AuthUseCase) Update(ctx context.Context, auth *domain.Auth) error {
	return a.repo.Update(ctx, auth)
}

func (a *AuthUseCase) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}

func (a *AuthUseCase) CheckAuth(ctx context.Context, uid string) (existed bool, err error) {
	auth, err := a.repo.GetById(ctx, uid)
	if err != nil {
		return false, domain.ErrAuthNotFound
	}
	if auth == nil {
		return false, nil
	}
	return true, nil
}

func (a *AuthUseCase) CheckAuthWithProvider(ctx context.Context, provider domain.Provider, token string) (existed bool, err error) {
	uid, _, err := provider.VerifyToken(ctx, token)
	if err != nil {
		return false, err
	}
	auth, err := a.repo.GetById(ctx, uid)
	if err != nil {
		return false, domain.ErrAuthNotFound
	}
	if auth == nil {
		return false, nil
	}
	return true, nil
}

func (a *AuthUseCase) SignUpWithProvider(ctx context.Context, provider domain.Provider, token string) error {
	uid, _, err := provider.VerifyToken(ctx, token)
	if err != nil {
		return err
	}
	// look for existing username
	found, err := a.repo.GetById(ctx, uid)
	if err == nil || found != nil {
		return domain.ErrUsernameExist
	}
	// create new auth
	auth := &domain.Auth{
		Id:       uid,
		Username: uid,
		RoleId:   a.defaultRoleId,
	}
	err = a.repo.Create(ctx, auth)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
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

func (a *AuthUseCase) SignUp(ctx context.Context, auth *domain.Auth) error {

	if auth.Id == "" {
		auth.Id = fmt.Sprintf("%d", time.Now().UnixMilli())
	}

	passwordPolicyErr := common.CheckPasswordPolicy(auth.Password, a.passwordPolicy)
	if passwordPolicyErr != nil {
		return passwordPolicyErr
	}
	// look for existing username
	found, err := a.repo.GetByUsername(ctx, auth.Username)
	if err == nil || found != nil {
		return domain.ErrUsernameExist
	}
	hashedPassword, err := common.GeneratePassword(auth.Password, a.hashCost)
	if err != nil {
		return err
	}
	auth.Hpassword = hashedPassword

	auth.RoleId = a.defaultRoleId

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
