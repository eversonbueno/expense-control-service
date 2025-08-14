package user

import (
	"context"
	"expense-control-service/internal/entity"
	userRepo "expense-control-service/internal/repositories/users"
)

type User interface {
	ListUsers(ctx context.Context) ([]*entity.User, error)
}

type user struct {
	userRepo userRepo.Users
}

func New(
	userRepo userRepo.Users,
) User {
	return &user{userRepo: userRepo}
}

func (u *user) ListUsers(ctx context.Context) ([]*entity.User, error)  {
	users, err := u.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
