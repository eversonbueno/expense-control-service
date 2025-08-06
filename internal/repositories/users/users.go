package users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eversonbueno/controle_gastos/internal/entity"
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
			&user.Email,
			&user.CPFCNPJ,
			&user.NomeCompleto,
			&user.Password,
			&user.SaldoUsuario,
			&user.TipoUsuario,
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
