package launch_category

import (
	"context"
	"expense-control-service/internal/entity"
	launchCategoryRepo "expense-control-service/internal/repositories/launch_category"
)

type LaunchCategory interface {
	ListLaunchCategorys(ctx context.Context) ([]*entity.LaunchCategory, error)
}

type launchCategory struct {
	launchCategoryRepo launchCategoryRepo.LaunchCategory
}

func New(
	launchCategoryRepo launchCategoryRepo.LaunchCategory,
) LaunchCategory {
	return &launchCategory{launchCategoryRepo: launchCategoryRepo}
}

func (u *launchCategory) ListLaunchCategorys(ctx context.Context) ([]*entity.LaunchCategory, error) {
	categories, err := u.launchCategoryRepo.ListLaunchCategorys(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
