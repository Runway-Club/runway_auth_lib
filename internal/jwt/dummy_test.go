package jwt_test

import (
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/Runway-Club/auth_lib/internal/jwt"
	"testing"
)

func TestDummyJwtGenerator(t *testing.T) {
	jwtDummy := jwt.NewDummyJwtGenerator("test", 3600000, "secret")
	t.Run("Generate and verify token", func(t *testing.T) {
		token, err := jwtDummy.GenerateToken(&domain.Auth{
			Id:       "user01",
			Username: "user01",
			RoleId:   "role01",
		}, map[string]interface{}{})
		if err != nil {
			t.Error(err)
		}
		auth, _, err := jwtDummy.VerifyToken(token)
		if err != nil {
			t.Error(err)
		}
		if auth == nil {
			t.Error("auth is nil")
		}
		if auth.Id != "user01" {
			t.Error("auth.Id is not user01")
		}
	})
}
