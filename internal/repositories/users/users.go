package users

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type Users interface {
	ListUserById(ctx context.Context, id int) (*entity.User, error)
}

type users struct {
	db *sql.DB
}

func New(db *sql.DB) Users  {
	return &users{db: db}
}

func (u *users) ListUserById(ctx context.Context, id int) (*entity.User, error)  {
	var user entity.User
	row := u.db.QueryRowContext(ctx, ListUserById, id)

	err := row.Scan(
		&user.Nome,
		&user.Sobrenome,
	)
	if err != nil {
			return nil, fmt.Errorf("erro ao scanear users: %v", err)
		}

	return &user, nil
}
