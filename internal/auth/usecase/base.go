package usecase

import (
	"context"
	"fmt"
	"github.com/Runway-Club/auth_lib/common"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/Runway-Club/auth_lib/utils"
	jwtlib "github.com/golang-jwt/jwt/v5"
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
	projectId      string
}

func (a *AuthUseCase) GetStaticUserList(ctx context.Context) (list *domain.StaticUserList, err error) {
	authList := make([]*domain.Auth, 0)
	for _, user := range a.repo.GetStaticUserMap(ctx) {
		authList = append(authList, user)
	}
	return &domain.StaticUserList{
		List: authList,
	}, nil
}

func (a *AuthUseCase) ChangePassword(ctx context.Context, uid, oldPassword, newPassword string) error {
	// get by username
	user, err := a.repo.GetById(ctx, uid)
	if err != nil {
		return domain.ErrAuthNotFound
	}

	// check password
	errPasswordPolicy := utils.CheckPasswordPolicy(newPassword, a.passwordPolicy)
	if errPasswordPolicy != nil {
		return errPasswordPolicy
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hpassword), []byte(oldPassword))
	if err != nil {
		return domain.ErrPasswordNotMatch
	}
	// hash password
	hashedPassword, err := utils.GeneratePassword(newPassword, a.hashCost)
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
	// Bearer process
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
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
	// prevent delete static user
	if _, ok := a.repo.GetStaticUserMap(ctx)[id]; ok {
		return domain.ErrPermissionDenied
	}
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

type CustomClaims struct {
	jwtlib.MapClaims
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func (a *AuthUseCase) customVerifyToken(ctx context.Context, token string) (uid string, claims map[string]interface{}, err error) {

	parsedClaims := &CustomClaims{}
	_, _, err = new(jwtlib.Parser).ParseUnverified(token, parsedClaims)

	log.Printf("parsedClaims: %v", parsedClaims)

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return "", nil, domain.ErrInvalidToken
	}

	// check expiration in seconds
	exp, err := parsedClaims.GetExpirationTime()
	if err != nil {
		return "", nil, domain.ErrInvalidToken
	}

	if time.Now().UTC().Unix() > exp.UTC().Unix() {
		return "", nil, domain.ErrExpiredToken

	}
	// check issuer
	issuer, err := parsedClaims.GetIssuer()
	if err != nil {
		return "", nil, domain.ErrInvalidToken
	}
	// Must be "https://securetoken.google.com/<projectId>", where <projectId> is the same project ID used for aud above.
	if issuer != fmt.Sprintf("https://securetoken.google.com/%s", a.projectId) {
		return "", nil, domain.ErrInvalidIssuer
	}

	uid = parsedClaims.Sub
	if uid == "" {
		return "", nil, domain.ErrInvalidToken
	}
	name := parsedClaims.Name
	// get email
	email := parsedClaims.Email
	if email == "" {
		return "", nil, domain.ErrInvalidToken
	}
	// get picture
	picture := parsedClaims.Picture
	return uid, map[string]interface{}{
		"name":    name,
		"email":   email,
		"picture": picture,
		"user_id": uid,
	}, nil

}

func (a *AuthUseCase) SignUpWithProvider(ctx context.Context, provider domain.Provider, token string) error {
	// bearer process
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	uid, _, err := provider.VerifyToken(ctx, token)
	if err != nil {
		// second chance for custom verify token
		uid, _, err = a.customVerifyToken(ctx, token)
		if err != nil {
			return err
		}
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
	// bearer process
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	uid, claims, err := provider.VerifyToken(ctx, token)
	if err != nil {
		// second chance for custom verify token
		uid, claims, err = a.customVerifyToken(ctx, token)
		if err != nil {
			return nil, err
		}
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

	passwordPolicyErr := utils.CheckPasswordPolicy(auth.Password, a.passwordPolicy)
	if passwordPolicyErr != nil {
		return passwordPolicyErr
	}
	// look for existing username
	found, err := a.repo.GetByUsername(ctx, auth.Username)
	if err == nil || found != nil {
		return domain.ErrUsernameExist
	}
	hashedPassword, err := utils.GeneratePassword(auth.Password, a.hashCost)
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
		projectId:      viper.GetString("runway_auth.projectid"),
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
