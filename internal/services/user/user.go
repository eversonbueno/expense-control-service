package user

import (
	"context"
	"expense-control-service/internal/entity"
	userRepo "expense-control-service/internal/repositories/users"
)

type User interface {
	ListUserById(ctx context.Context, id int) (*entity.User, error)
}

type user struct {
	userRepo userRepo.Users
}

func New(
	userRepo userRepo.Users,
) User {
	return &user{userRepo: userRepo}
}

func (u *user) ListUserById(ctx context.Context, id int) (*entity.User, error)  {
	users, err := u.userRepo.ListUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return users, nil
}
