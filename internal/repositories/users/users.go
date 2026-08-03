package users

import (
	"context"
	"database/sql"
	"errors"
	"expense-control-service/internal/entity"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// ErrEmailDuplicado é retornado por Create quando a constraint UNIQUE(email) do banco
// rejeita a inserção — cobre a corrida entre o FindByEmail do usecase e este INSERT.
var ErrEmailDuplicado = errors.New("email duplicado")

const mysqlErrDuplicateEntry = 1062

type Users interface {
	ListUserById(ctx context.Context, id int) (*entity.User, error)
	ListUsers(ctx context.Context) ([]*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
}

type users struct {
	db *sql.DB
}

func New(db *sql.DB) Users {
	return &users{db: db}
}

func scanUser(row *sql.Row, user *entity.User) error {
	return row.Scan(
		&user.ID,
		&user.Nome,
		&user.Email,
		&user.SenhaHash,
		&user.GrupoFamiliarID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func (u *users) ListUserById(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	row := u.db.QueryRowContext(ctx, ListUserById, id)

	if err := scanUser(row, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao scanear users: %v", err)
	}

	return &user, nil
}

func (u *users) ListUsers(ctx context.Context) ([]*entity.User, error) {
	rows, err := u.db.QueryContext(ctx, ListUsers)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %v", err)
	}
	defer rows.Close()

	var result []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID,
			&user.Nome,
			&user.Email,
			&user.SenhaHash,
			&user.GrupoFamiliarID,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao scanear usuário: %v", err)
		}
		result = append(result, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar usuários: %v", err)
	}

	return result, nil
}

// FindByEmail retorna nil, nil quando nenhum usuário é encontrado com o e-mail informado.
func (u *users) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	row := u.db.QueryRowContext(ctx, FindByEmail, email)

	if err := scanUser(row, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar usuário por email: %v", err)
	}

	return &user, nil
}

func (u *users) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	result, err := u.db.ExecContext(ctx, CreateUser, user.Nome, user.Email, user.SenhaHash, user.GrupoFamiliarID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlErrDuplicateEntry {
			return nil, ErrEmailDuplicado
		}
		return nil, fmt.Errorf("erro ao criar usuário: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter id do usuário criado: %v", err)
	}

	return u.ListUserById(ctx, int(id))
}
