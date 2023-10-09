package domain

import (
	"context"
	"errors"
	"github.com/Runway-Club/auth_lib/common"
	"gorm.io/gorm"
)

type ACI struct {
	gorm.Model
	Id       string `json:"id" gorm:"uniqueIndex" yaml:"id"`
	Resource string `json:"resource" yaml:"resource"`
	Payload  string `json:"payload" yaml:"payload"`
	RoleId   string `json:"role_id" yaml:"roleId"`
	UserId   string `json:"user_id" yaml:"userId"`
}

type ACIRepository interface {
	Create(ctx context.Context, aci *ACI) error
	GetById(ctx context.Context, id string) (*ACI, error)
	GetByResource(ctx context.Context, resource string) ([]*ACI, error)
	GetByRoleId(ctx context.Context, roleId string) ([]*ACI, error)
	GetByPayload(ctx context.Context, payload string) ([]*ACI, error)
	GetByUserId(ctx context.Context, userId string) ([]*ACI, error)
	CheckByRoleId(ctx context.Context, roleId string, resource string, payload string) (bool, error)
	CheckByUserId(ctx context.Context, userId string, resource string, payload string) (bool, error)
	List(ctx context.Context, query *common.QueryOpts) (*common.ListResult[*ACI], error)
	Update(ctx context.Context, aci *ACI) error
	Delete(ctx context.Context, id string) error
	GetResourcesByUserIdAndPayload(ctx context.Context, userId string, payload string) ([]*ACI, error)
}

type ACIUseCase interface {
	Create(ctx context.Context, aci *ACI) error
	GetById(ctx context.Context, id string) (*ACI, error)
	GetByResource(ctx context.Context, resource string) ([]*ACI, error)
	GetByRoleId(ctx context.Context, roleId string) ([]*ACI, error)
	GetByPayload(ctx context.Context, payload string) ([]*ACI, error)
	GetByUserId(ctx context.Context, userId string) ([]*ACI, error)
	GetResourcesByUserIdAndPayload(ctx context.Context, userId string, payload string) ([]*ACI, error)
	List(ctx context.Context, query *common.QueryOpts) (*common.ListResult[*ACI], error)
	Update(ctx context.Context, aci *ACI) error
	Delete(ctx context.Context, id string) error
}

var (
	ErrACINotFound      = errors.New("aci not found")
	ErrPermissionDenied = errors.New("permission denied")
)
