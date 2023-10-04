package repo

import (
	"context"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ACIRepository struct {
	db *gorm.DB
}

type defaultACL struct {
	ACL []domain.ACI `json:"acl" yaml:"acl"`
}

func NewACIRepository(dialector gorm.Dialector) *ACIRepository {
	db, err := gorm.Open(dialector)
	if err != nil {
		panic(err)
	}
	// migrate schema
	err = db.AutoMigrate(&domain.ACI{})
	// init default aci
	acl := defaultACL{}
	viper.UnmarshalKey("runway_auth", &acl)
	for _, item := range acl.ACL {
		db.Create(&item)
	}
	return &ACIRepository{
		db: db,
	}
}

func (a *ACIRepository) Create(ctx context.Context, aci *domain.ACI) error {
	tx := a.db.WithContext(ctx).Create(aci)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *ACIRepository) GetById(ctx context.Context, id string) (*domain.ACI, error) {
	found := &domain.ACI{}
	tx := a.db.WithContext(ctx).Where("id = ?", id).First(found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *ACIRepository) GetByResource(ctx context.Context, resource string) ([]*domain.ACI, error) {
	var found []*domain.ACI
	tx := a.db.WithContext(ctx).Where("resource = ?", resource).Find(&found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *ACIRepository) GetByRoleId(ctx context.Context, roleId string) ([]*domain.ACI, error) {
	var found []*domain.ACI
	tx := a.db.WithContext(ctx).Where("role_id = ?", roleId).Find(&found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *ACIRepository) GetByPayload(ctx context.Context, payload string) ([]*domain.ACI, error) {
	var found []*domain.ACI
	tx := a.db.WithContext(ctx).Where("payload = ?", payload).Find(&found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *ACIRepository) GetByUserId(ctx context.Context, userId string) ([]*domain.ACI, error) {
	var found []*domain.ACI
	tx := a.db.WithContext(ctx).Where("user_id = ?", userId).Find(&found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *ACIRepository) CheckByRoleId(ctx context.Context, roleId string, resource string, payload string) (bool, error) {
	found := &domain.ACI{}
	tx := a.db.WithContext(ctx).Where("role_id = ? AND resource = ? AND payload = ?", roleId, resource, payload).First(found)
	if tx.Error != nil {
		return false, tx.Error
	}
	return true, nil
}

func (a *ACIRepository) CheckByUserId(ctx context.Context, userId string, resource string, payload string) (bool, error) {
	found := &domain.ACI{}
	tx := a.db.WithContext(ctx).Where("user_id = ? AND resource = ? AND payload = ?", userId, resource, payload).First(found)
	if tx.Error != nil {
		return false, tx.Error
	}
	return true, nil
}
