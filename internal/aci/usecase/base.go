package usecase

import (
	"context"
	"fmt"
	"github.com/Runway-Club/auth_lib/common"
	"github.com/Runway-Club/auth_lib/domain"
	"time"
)

type ACIUseCase struct {
	aciRepo domain.ACIRepository
}

func (a *ACIUseCase) GetResourcesByUserIdAndPayload(ctx context.Context, userId string, payload string) ([]*domain.ACI, error) {
	return a.aciRepo.GetResourcesByUserIdAndPayload(ctx, userId, payload)
}

func (a *ACIUseCase) Update(ctx context.Context, aci *domain.ACI) error {
	if aci.Id == "" {
		return domain.ErrInvalidACI
	}
	// check if aci exist
	foundACI, err := a.aciRepo.GetById(ctx, aci.Id)
	if err != nil || foundACI == nil {
		return domain.ErrACINotFound
	}

	return a.aciRepo.Update(ctx, aci)
}

func (a *ACIUseCase) Delete(ctx context.Context, id string) error {
	return a.aciRepo.Delete(ctx, id)
}

func (a *ACIUseCase) List(ctx context.Context, query *common.QueryOpts) (*common.ListResult[*domain.ACI], error) {
	return a.aciRepo.List(ctx, query)
}

func (a *ACIUseCase) Create(ctx context.Context, aci *domain.ACI) error {
	if aci.Resource == "" {
		return domain.ErrInvalidACI
	}
	// create random id if aci.id is empty
	if aci.Id == "" {
		aci.Id = fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	return a.aciRepo.Create(ctx, aci)
}

func (a *ACIUseCase) GetById(ctx context.Context, id string) (*domain.ACI, error) {
	return a.aciRepo.GetById(ctx, id)
}

func (a *ACIUseCase) GetByResource(ctx context.Context, resource string) ([]*domain.ACI, error) {
	return a.aciRepo.GetByResource(ctx, resource)
}

func (a *ACIUseCase) GetByRoleId(ctx context.Context, roleId string) ([]*domain.ACI, error) {
	return a.aciRepo.GetByRoleId(ctx, roleId)
}

func (a *ACIUseCase) GetByPayload(ctx context.Context, payload string) ([]*domain.ACI, error) {
	return a.aciRepo.GetByPayload(ctx, payload)
}

func (a *ACIUseCase) GetByUserId(ctx context.Context, userId string) ([]*domain.ACI, error) {
	return a.aciRepo.GetByUserId(ctx, userId)
}

func NewACIUseCase(aciRepo domain.ACIRepository) *ACIUseCase {
	return &ACIUseCase{
		aciRepo: aciRepo,
	}
}
