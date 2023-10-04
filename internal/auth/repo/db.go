package repo

import (
	"context"
	"gorm.io/gorm"
	"runwayclub.dev/auth/domain"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(dialector gorm.Dialector) *AuthRepository {
	db, err := gorm.Open(dialector)
	if err != nil {
		panic(err)
	}
	// migrate schema
	err = db.AutoMigrate(&domain.Auth{})
	if err != nil {
		panic(err)
	}
	return &AuthRepository{db: db}
}

func (a *AuthRepository) Create(ctx context.Context, auth *domain.Auth) error {
	tx := a.db.WithContext(ctx).Create(auth)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *AuthRepository) GetById(ctx context.Context, id string) (*domain.Auth, error) {
	found := &domain.Auth{}
	tx := a.db.WithContext(ctx).Where("id = ?", id).First(found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *AuthRepository) GetByUsernameAndHpassword(ctx context.Context, username, hpassword string) (*domain.Auth, error) {
	found := &domain.Auth{}
	tx := a.db.WithContext(ctx).Where("username = ? AND hpassword = ?", username, hpassword).First(found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *AuthRepository) GetByUsername(ctx context.Context, username string) (*domain.Auth, error) {
	found := &domain.Auth{}
	tx := a.db.WithContext(ctx).Where("username = ?", username).First(found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *AuthRepository) Update(ctx context.Context, auth *domain.Auth) error {
	tx := a.db.WithContext(ctx).Save(auth)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *AuthRepository) Delete(ctx context.Context, id string) error {
	tx := a.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Auth{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
