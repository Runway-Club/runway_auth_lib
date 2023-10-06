package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/Runway-Club/auth_lib/internal/auth/repo"
	"github.com/Runway-Club/auth_lib/internal/auth/usecase"
	"github.com/Runway-Club/auth_lib/internal/jwt"
	"github.com/Runway-Club/auth_lib/internal/providers"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"testing"
)

func TestAuthUseCase(t *testing.T) {
	viper.SetConfigFile("../../../configs/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	authRepo := repo.NewAuthRepository(sqlite.Open(":memory:"))
	authUseCase := usecase.NewAuthUseCase(authRepo, jwt.NewJwtGenerator())

	dummyJwtGenerator := jwt.NewDummyJwtGenerator("test.com", 3600000, "secret")
	provider := providers.NewDummyProvider(dummyJwtGenerator)

	t.Run("sign up", func(t *testing.T) {
		err := authUseCase.SignUp(context.Background(), &domain.Auth{
			Id:       "1",
			Username: "test",
			Password: "test12345678",
		})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("sign up with existing username", func(t *testing.T) {
		err := authUseCase.SignUp(context.Background(), &domain.Auth{
			Id:       "2",
			Username: "test2",
			Password: "test12345678",
		})
		if errors.Is(err, domain.ErrUsernameExist) {
			t.Error("expected error username already exist")
		}
	})
	t.Run("sign up with invalid password (Level 2)", func(t *testing.T) {
		viper.Set("runway_auth.password.policy", "level2")
		authUseCase = usecase.NewAuthUseCase(authRepo, jwt.NewJwtGenerator())
		err := authUseCase.SignUp(context.Background(), &domain.Auth{
			Id:       "3",
			Username: "test3",
			Password: "thisisinvalidpassword",
		})
		if !errors.Is(err, domain.ErrInvalidPassword) {
			t.Error("expected error invalid password")
		}
		err = authUseCase.SignUp(context.Background(), &domain.Auth{
			Id:       "3",
			Username: "test4",
			Password: "thisisvalidpassword1",
		})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("sign in", func(t *testing.T) {
		token, err := authUseCase.SignIn(context.Background(), "test", "test12345678")
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			t.Error("token is empty")
		}
		fmt.Println(token)
	})

	t.Run("Sign up with provider", func(t *testing.T) {
		// generate token from dummy provider
		token, err := dummyJwtGenerator.GenerateToken(&domain.Auth{
			Id:     "0002",
			RoleId: "admin",
		}, map[string]interface{}{})
		if err != nil {
			t.Error(err)
		}
		err = authUseCase.SignUpWithProvider(context.Background(), provider, token)
		if err != nil {
			t.Error(err)
		}
		// get new user by username
		user, err := authRepo.GetByUsername(context.Background(), "0002")
		if err != nil {
			t.Error(err)
		}
		if user == nil {
			t.Error("user is empty")
		}
		if user.Id != "0002" {
			t.Error("user id is not 0002")
		}
	})

	t.Run("Sign in with provider", func(t *testing.T) {
		token, err := dummyJwtGenerator.GenerateToken(&domain.Auth{
			Id:     "0002",
			RoleId: "admin",
		}, map[string]interface{}{})
		if err != nil {
			t.Error(err)
		}
		auth, err := authUseCase.SignInWithProvider(context.Background(), provider, token)
		if err != nil {
			t.Error(err)
		}
		if auth == nil {
			t.Error("auth is empty")
		}
		if auth.Id != "0002" {
			t.Error("auth id is not 0002")
		}
	})

}
