package runway_auth_test

import (
	"context"
	auth "github.com/Runway-Club/auth_lib"
	"github.com/Runway-Club/auth_lib/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestAuth(t *testing.T) {
	auth.Initialize("configs/dev.yaml", func() gorm.Dialector {
		return sqlite.Open("file::memory:?cache=shared")
	}, func() gorm.Dialector {
		return sqlite.Open("file::memory:?cache=shared")
	}, true)
	t.Run("sign up", func(t *testing.T) {
		err := auth.SignUp(context.Background(), &domain.Auth{
			Id:       "test",
			Username: "user01",
			Password: "Strong123456",
		})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("sign in", func(t *testing.T) {
		token, err := auth.SignIn(context.Background(), "user01", "Strong123456")
		if err != nil {
			t.Error(err)
		}
		if token == nil {
			t.Error("token is nil")
		}
		if token.Id != "test" {
			t.Errorf("expected id 1, got %s", token.Id)
		}
		if token.UserId != "test" {
			t.Errorf("expected user id 1, got %s", token.UserId)
		}
		if token.RoleId != "default" {
			t.Errorf("expected role id default, got %s", token.RoleId)
		}
		if token.Jwt == "" {
			t.Error("jwt is empty")
		}
	})

	t.Run("Create aci", func(t *testing.T) {
		err := auth.GetACIUseCase().Create(context.Background(), &domain.ACI{
			Id:       "100",
			Resource: "v1/course.GET",
			RoleId:   "default",
		})
		if err != nil {
			t.Error(err)
		}
		err = auth.GetACIUseCase().Create(context.Background(), &domain.ACI{
			Id:       "101",
			Resource: "v1/course.POST",
			RoleId:   "admin",
		})
		if err != nil {
			t.Error(err)
		}
		err = auth.GetACIUseCase().Create(context.Background(), &domain.ACI{
			Id:       "102",
			Resource: "v1/course.PUT",
			RoleId:   "admin",
			UserId:   "1",
		})
	})
	t.Run("verify token and check perm", func(t *testing.T) {
		token, err := auth.SignIn(context.Background(), "user01", "Strong123456")
		if err != nil {
			t.Error(err)
		}
		result := auth.VerifyTokenAndPerm(context.Background(), token.Jwt, "v1/course.GET", "")
		if result != nil {
			t.Error("expected true, got false")
		}

		result = auth.VerifyTokenAndPerm(context.Background(), token.Jwt, "v1/course.POST", "demo")
		if result != nil {
			t.Error("expected true, got false")
		}

		result = auth.VerifyTokenAndPerm(context.Background(), token.Jwt, "v1/course.PUT", "")
		if result != nil {
			t.Error("expected true, got false")
		}

	})
	t.Run("bypass with admin user", func(t *testing.T) {
		token, err := auth.SignIn(context.Background(), "admin", "Adminpassword@123")
		if err != nil {
			t.Error(err)
		}
		result := auth.VerifyTokenAndPerm(context.Background(), token.Jwt, "v1/payment.POST", "")
		if result != nil {
			t.Error("expected true, got false")
		}
	})
}
