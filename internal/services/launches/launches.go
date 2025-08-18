package launches

import (
	"context"
	"expense-control-service/internal/entity"
	launchesRepo "expense-control-service/internal/repositories/launches"
)

type Launches interface {
	ListLaunches(ctx context.Context) ([]*entity.Launches, error)
}

type launches struct {
	launchesRepo launchesRepo.Launches
}

func New(
	launchesRepo launchesRepo.Launches,
) Launches  {
	return &launches{launchesRepo: launchesRepo}
}

func (l *launches) ListLaunches(ctx context.Context) ([]*entity.Launches, error)  {
	launches,err := l.launchesRepo.ListLaunches(ctx)
	if err != nil {
		return nil, err
	}

	return launches, nil
}
