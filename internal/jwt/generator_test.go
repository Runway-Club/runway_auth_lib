package jwt_test

import (
	"fmt"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/Runway-Club/auth_lib/internal/jwt"
	"github.com/spf13/viper"
	"testing"
)

func TestJwtGenerator(t *testing.T) {
	viper.Set("runway_auth.jwt.secret", "this-is-a-secret")
	viper.Set("runway_auth.jwt.exp", 1000)
	viper.Set("runway_auth.jwt.issuer", "runwayclub")
	generator := jwt.NewJwtGenerator()
	t.Run("generate token", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "test",
			"yob":  1999,
		}
		token, err := generator.GenerateToken(&domain.Auth{
			Id:       "1",
			Username: "test",
			RoleId:   "1",
		}, data)
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("token is empty")
		}
		fmt.Println(token)
		_, parsedData, err := generator.VerifyToken(token)
		if err != nil {
			t.Error(err)
		}
		if parsedData["name"] != "test" {
			t.Errorf("expected name test, got %s", parsedData["name"])
		}
	})
}
