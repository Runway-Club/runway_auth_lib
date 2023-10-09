package domain

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	Id        string `json:"id" gorm:"uniqueIndex" mapstructure:"id"`
	Username  string `json:"username" gorm:"uniqueIndex"`
	Password  string `json:"password" gorm:"-"`
	Hpassword string `json:"hpassword"`
	RoleId    string `json:"role_id" mapstructure:"role_id"`
}

type Token struct {
	Jwt    string `json:"jwt"`
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	RoleId string `json:"role_id"`
}

type StaticUserList struct {
	List []*Auth `json:"static_users" yaml:"static_users" mapstructure:"static_users"`
}

type AuthRepository interface {
	Create(ctx context.Context, auth *Auth) error
	GetById(ctx context.Context, id string) (*Auth, error)
	GetStaticUserMap(ctx context.Context) map[string]*Auth
	GetByUsername(ctx context.Context, username string) (*Auth, error)
	GetByUsernameAndHpassword(ctx context.Context, username, hpassword string) (*Auth, error)
	Update(ctx context.Context, auth *Auth) error
	Delete(ctx context.Context, id string) error
}

type AuthUseCase interface {
	SignUp(ctx context.Context, auth *Auth) error
	SignUpWithProvider(ctx context.Context, provider Provider, token string) error
	SignIn(ctx context.Context, username, password string) (token *Token, err error)
	SignInWithProvider(ctx context.Context, provider Provider, token string) (genToken *Token, err error)
	CheckAuth(ctx context.Context, uid string) (existed bool, err error)
	CheckAuthWithProvider(ctx context.Context, provider Provider, token string) (existed bool, err error)
}

var (
	ErrAuthNotFound     = errors.New("auth not found")
	ErrUsernameExist    = errors.New("username already exist")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrPasswordNotMatch = errors.New("password not match")
	ErrInternal         = errors.New("internal error")
)
