package repo

import (
	"context"
	"github.com/Runway-Club/auth_lib/common"
	"github.com/Runway-Club/auth_lib/domain"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"math"
)

type AuthRepository struct {
	db          *gorm.DB
	StaticUsers *domain.StaticUserList
	UserIdMap   map[string]*domain.Auth
}

func (a *AuthRepository) List(ctx context.Context, opt *common.QueryOpts) (*common.ListResult[*domain.Auth], error) {
	var auths []*domain.Auth
	tx := a.db.WithContext(ctx)
	offset := opt.Size * (opt.Page - 1)
	if opt != nil {
		tx = tx.Limit(opt.Size).Offset(offset)
	}
	tx = tx.Find(&auths)
	if tx.Error != nil {
		return nil, tx.Error
	}
	count := int64(0)
	tx = a.db.WithContext(ctx).Model(&domain.Auth{}).Count(&count)
	if tx.Error != nil {
		return nil, tx.Error
	}
	endPage := int(math.Floor(float64(count) / float64(opt.Size)))
	return &common.ListResult[*domain.Auth]{
		Data:    auths,
		EndPage: endPage,
	}, nil
}

func NewAuthRepository(dialector gorm.Dialector) *AuthRepository {
	// load static users
	staticUsers := &domain.StaticUserList{
		List: make([]*domain.Auth, 0),
	}
	userIdMap := make(map[string]*domain.Auth)
	err := viper.UnmarshalKey("runway_auth.static_users", &staticUsers.List)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(dialector)
	if err != nil {
		panic(err)
	}
	// migrate schema
	err = db.AutoMigrate(&domain.Auth{})
	if err != nil {
		panic(err)
	}
	// create static users, omit error because it's okay if it's already exist
	for _, user := range staticUsers.List {
		userIdMap[user.Id] = user
	}
	return &AuthRepository{db: db, StaticUsers: staticUsers, UserIdMap: userIdMap}
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

func (a *AuthRepository) GetStaticUserMap(ctx context.Context) map[string]*domain.Auth {
	return a.UserIdMap
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
	if tx.RowsAffected == 0 {
		return domain.ErrAuthNotFound
	}
	return nil
}
