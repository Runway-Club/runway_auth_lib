package usecase

import (
	"context"
	"github.com/Runway-Club/auth_lib/domain"
)

type ACIUseCase struct {
	aciRepo domain.ACIRepository
}

func (a *ACIUseCase) Create(ctx context.Context, aci *domain.ACI) error {
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
