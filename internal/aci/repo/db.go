package repo

import (
	"context"
	"github.com/Runway-Club/auth_lib/common"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"math"
)

type ACIRepository struct {
	db *gorm.DB
}

func (a *ACIRepository) GetResourcesByUserIdAndPayload(ctx context.Context, userId string, payload string) ([]*domain.ACI, error) {
	var found []*domain.ACI
	tx := a.db.WithContext(ctx).Where("user_id = ? AND payload = ?", userId, payload).Find(&found)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return found, nil
}

func (a *ACIRepository) Update(ctx context.Context, aci *domain.ACI) error {
	tx := a.db.WithContext(ctx).Save(aci)
	return tx.Error
}

func (a *ACIRepository) Delete(ctx context.Context, id string) error {
	tx := *a.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.ACI{})
	return tx.Error
}

func (a *ACIRepository) List(ctx context.Context, query *common.QueryOpts) (*common.ListResult[*domain.ACI], error) {
	acl := make([]*domain.ACI, 0)
	offset := query.Page * query.Size
	tx := a.db.WithContext(ctx).Offset(offset).Limit(query.Size).Find(acl)
	if tx.Error != nil {
		return nil, tx.Error
	}
	count := int64(0)
	// count all row
	a.db.WithContext(ctx).Count(&count)
	numOfPage := int(math.Ceil(float64(count) / float64(query.Size)))

	return &common.ListResult[*domain.ACI]{
		Data:    acl,
		EndPage: numOfPage,
	}, nil
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
	return tx.RowsAffected == 0, tx.Error
}

func (a *ACIRepository) CheckByUserId(ctx context.Context, userId string, resource string, payload string) (bool, error) {
	found := &domain.ACI{}
	tx := a.db.WithContext(ctx).Where("user_id = ? AND resource = ? AND payload = ?", userId, resource, payload).First(found)
	return tx.RowsAffected == 0, tx.Error
}
