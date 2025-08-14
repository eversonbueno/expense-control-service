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
			&user.NomeCompleto,
			&user.CPFCNPJ,
			&user.Email,
			&user.Password,
			&user.TipoUsuario,
			&user.SaldoUsuario,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear usuarios: %v", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar linhas: %v", err)
	}

	return users, nil
}
