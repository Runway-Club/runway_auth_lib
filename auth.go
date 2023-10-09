package runway_auth

import (
	"context"
	"github.com/Runway-Club/auth_lib/domain"
	aciRepoPkg "github.com/Runway-Club/auth_lib/internal/aci/repo"
	aciUseCasePkg "github.com/Runway-Club/auth_lib/internal/aci/usecase"
	authRepoPkg "github.com/Runway-Club/auth_lib/internal/auth/repo"
	authUseCasePkg "github.com/Runway-Club/auth_lib/internal/auth/usecase"
	jwtPkg "github.com/Runway-Club/auth_lib/internal/jwt"
	providerPkg "github.com/Runway-Club/auth_lib/internal/providers"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	authRepo     domain.AuthRepository
	aciRepo      domain.ACIRepository
	authUseCase  domain.AuthUseCase
	aciUseCase   domain.ACIUseCase
	jwtGenerator domain.JwtGenerator
	provider     domain.Provider
)

type InitDialetor func() gorm.Dialector

func Initialize(configFileName string, authDialector InitDialetor, aciDialector InitDialetor, setConfig bool) {
	if setConfig {
		viper.SetConfigFile(configFileName)
		err := viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
	}
	jwtGenerator = jwtPkg.NewJwtGenerator()
	authRepo = authRepoPkg.NewAuthRepository(authDialector())
	aciRepo = aciRepoPkg.NewACIRepository(aciDialector())
	authUseCase = authUseCasePkg.NewAuthUseCase(authRepo, jwtGenerator)
	aciUseCase = aciUseCasePkg.NewACIUseCase(aciRepo)
}

func SignUp(ctx context.Context, auth *domain.Auth) error {
	return authUseCase.SignUp(ctx, auth)
}

func SignUpWithProvider(ctx context.Context, token string) error {
	if provider == nil {
		panic("provider not initialized")
	}
	return authUseCase.SignUpWithProvider(ctx, provider, token)
}

func SignInWithProvider(ctx context.Context, token string) (genToken *domain.Token, err error) {
	if provider == nil {
		panic("provider not initialized")
	}
	return authUseCase.SignInWithProvider(ctx, provider, token)
}

func SignIn(ctx context.Context, username, password string) (token *domain.Token, err error) {
	return authUseCase.SignIn(ctx, username, password)
}

func InitGoogleProvider(ctx context.Context, firebaseAdminConfigName string) {
	provider = providerPkg.NewGoogleProvider(ctx, firebaseAdminConfigName)
}

func GetAuthUseCase() domain.AuthUseCase {
	return authUseCase
}

func GetACIUseCase() domain.ACIUseCase {
	return aciUseCase
}

func VerifyTokenAndPerm(ctx context.Context, token, resource, payload string) error {
	auth, _, err := jwtGenerator.VerifyToken(token)
	// bypass if auth is static
	if len(authRepo.GetStaticUserMap(ctx)) > 0 {
		if _, ok := authRepo.GetStaticUserMap(ctx)[auth.Id]; ok {
			return nil
		}
	}

	if err != nil {
		return err
	}
	result, err := aciRepo.CheckByUserId(ctx, auth.Id, resource, payload)
	if result {
		return nil
	}
	result, err = aciRepo.CheckByRoleId(ctx, auth.RoleId, resource, payload)
	if result {
		return nil
	}
	return domain.ErrPermissionDenied
}

func CheckAuthWithProvider(ctx context.Context, token string) (bool, error) {
	if provider == nil {
		panic("provider not initialized")
	}
	return authUseCase.CheckAuthWithProvider(ctx, provider, token)
}

func VerifyToken(ctx context.Context, token string) (auth *domain.Auth, err error) {
	return authUseCase.Verify(ctx, token)
}
