package users

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type Users interface {
	ListUsers(ctx context.Context) ([]*entity.User, error)
}

type users struct {
	db *sql.DB
}

func New(db *sql.DB) Users  {
	return &users{db: db}
}

func (u *users) ListUsers(ctx context.Context) ([]*entity.User, error)  {
	rows, err := u.db.QueryContext(ctx, ListUsers)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %v", err)
	}

	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Nome,
			&user.Sobrenome,
			&user.Usuario,
			&user.Senha,
			&user.Saldo,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear users: %v", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar linhas: %v", err)
	}

	return users, nil
}
