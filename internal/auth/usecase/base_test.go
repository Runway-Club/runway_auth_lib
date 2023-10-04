package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/Runway-Club/auth_lib/internal/auth/repo"
	"github.com/Runway-Club/auth_lib/internal/auth/usecase"
	"github.com/Runway-Club/auth_lib/internal/jwt"
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
}
