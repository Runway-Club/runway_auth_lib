package repo_test

import (
	"context"
	"gorm.io/driver/sqlite"
	"runwayclub.dev/auth/domain"
	"runwayclub.dev/auth/internal/auth/repo"
	"testing"
)

func TestNewAuthRepository(t *testing.T) {
	// gorm in memory db
	sqliteDialector := sqlite.Open(":memory:")
	dbRepo := repo.NewAuthRepository(sqliteDialector)
	t.Run("create auth", func(t *testing.T) {
		newUser := &domain.Auth{
			Id:        "1",
			Username:  "test",
			Password:  "test",
			Hpassword: "test",
		}
		err := dbRepo.Create(context.Background(), newUser)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("get auth by id", func(t *testing.T) {
		found, err := dbRepo.GetById(context.Background(), "1")
		if err != nil {
			t.Error(err)
		}
		if found.Username != "test" {
			t.Errorf("expected username test, got %s", found.Username)
		}
	})
	t.Run("get auth by username and hpassword", func(t *testing.T) {
		found, err := dbRepo.GetByUsernameAndHpassword(context.Background(), "test", "test")
		if err != nil {
			t.Error(err)
		}
		if found.Username != "test" {
			t.Errorf("expected username test, got %s", found.Username)
		}
	})
	t.Run("update auth", func(t *testing.T) {
		found, err := dbRepo.GetById(context.Background(), "1")
		if err != nil {
			t.Error(err)
		}
		found.Username = "test2"
		err = dbRepo.Update(context.Background(), found)
		if err != nil {
			t.Error(err)
		}
		found, err = dbRepo.GetById(context.Background(), "1")
		if err != nil {
			t.Error(err)
		}
		if found.Username != "test2" {
			t.Errorf("expected username test2, got %s", found.Username)
		}
	})
	t.Run("delete auth", func(t *testing.T) {

		dbRepo.Create(context.Background(), &domain.Auth{
			Id:       "2",
			Username: "test4",
			Password: "test4",
		})

		err := dbRepo.Delete(context.Background(), "2")
		if err != nil {
			t.Error(err)
		}
		_, err = dbRepo.GetById(context.Background(), "2")
		if err == nil {
			t.Error("expected error, got nil")
		}

	})
}
