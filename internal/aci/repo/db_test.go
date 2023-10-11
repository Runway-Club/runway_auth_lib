package repo_test

import (
	"context"
	"github.com/Runway-Club/auth_lib/common"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/Runway-Club/auth_lib/internal/aci/repo"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"testing"
)

func TestACIRepository(t *testing.T) {
	viper.SetConfigFile("../../../configs/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	sqliteDialector := sqlite.Open(":memory:")
	aciRepo := repo.NewACIRepository(sqliteDialector)
	t.Run("create aci", func(t *testing.T) {
		err := aciRepo.Create(nil, &domain.ACI{
			Id:       "100",
			Resource: "test",
			Payload:  "test",
		})
		if err != nil {
			t.Error(err)
		}
		result, err := aciRepo.List(context.Background(), &common.QueryOpts{
			Page: 1,
			Size: 10,
		})
		if err != nil {
			t.Error(err)
		}
		if len(result.Data) != 1 {
			t.Error("invalid result")
		}
	})
}
