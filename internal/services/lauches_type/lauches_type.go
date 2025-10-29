package lauches_type

import (
	"context"
	"expense-control-service/internal/entity"
	lauchTypesRepo "expense-control-service/internal/repositories/launch_type"
)

type LauchType interface {
	ListLauchTypesById(ctx context.Context, id int) (*entity.LauchType, error)
}

type lauchType struct {
	lauchTypeRepo lauchTypesRepo.LauchType
}

func New(
	lauchTypeRepo lauchTypesRepo.LauchType,
) LauchType {
	return &lauchType{lauchTypeRepo: lauchTypeRepo}
}

func (u *lauchType) ListLauchTypesById(ctx context.Context, id int) (*entity.LauchType, error)  {
	paymentMethod, err := u.lauchTypeRepo.ListLauchTypesById(ctx, id)
	if err != nil {
		return nil, err
	}

	return paymentMethod, nil
}
