package user

import (
	"context"
	"github.com/eversonbueno/controle_gastos/internal/entity"
	userRepo "github.com/eversonbueno/controle_gastos/internal/repositories/users"
)

type User interface {

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
	users, err := u.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
