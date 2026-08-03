package convite

import (
	"context"
	"database/sql"
	"errors"
	"expense-control-service/internal/entity"
	"fmt"
	"time"
)

// ErrJaUtilizado é retornado por MarcarUtilizado quando o convite já havia sido
// consumido por outra requisição concorrente (a UPDATE condicional não afetou linhas).
var ErrJaUtilizado = errors.New("convite já utilizado")

type Convite interface {
	Create(ctx context.Context, convite *entity.Convite) (*entity.Convite, error)
	FindByCodigo(ctx context.Context, codigo string) (*entity.Convite, error)
	MarcarUtilizado(ctx context.Context, id int, utilizadoEm time.Time) error
}

type convite struct {
	db *sql.DB
}

func New(db *sql.DB) Convite {
	return &convite{db: db}
}

func scanConvite(row *sql.Row, c *entity.Convite) error {
	var utilizadoEm sql.NullTime
	if err := row.Scan(&c.ID, &c.Codigo, &c.GrupoFamiliarID, &c.CriadoEm, &c.ExpiraEm, &utilizadoEm); err != nil {
		return err
	}
	if utilizadoEm.Valid {
		c.UtilizadoEm = &utilizadoEm.Time
	}
	return nil
}

func (r *convite) Create(ctx context.Context, c *entity.Convite) (*entity.Convite, error) {
	result, err := r.db.ExecContext(ctx, CreateConvite, c.Codigo, c.GrupoFamiliarID, c.ExpiraEm)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar convite: %v", err)
	}

	if _, err := result.LastInsertId(); err != nil {
		return nil, fmt.Errorf("erro ao obter id do convite criado: %v", err)
	}

	return r.FindByCodigo(ctx, c.Codigo)
}

// FindByCodigo retorna nil, nil quando nenhum convite é encontrado com o código informado.
func (r *convite) FindByCodigo(ctx context.Context, codigo string) (*entity.Convite, error) {
	var c entity.Convite
	row := r.db.QueryRowContext(ctx, FindByCodigo, codigo)

	if err := scanConvite(row, &c); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar convite por código: %v", err)
	}

	return &c, nil
}

// MarcarUtilizado marca o convite como consumido de forma atômica (UPDATE condicional
// a utilizado_em IS NULL). Retorna ErrJaUtilizado se outra requisição concorrente já
// tiver consumido o mesmo convite entre a checagem em FindByCodigo e esta chamada.
func (r *convite) MarcarUtilizado(ctx context.Context, id int, utilizadoEm time.Time) error {
	result, err := r.db.ExecContext(ctx, MarcarUtilizado, utilizadoEm, id)
	if err != nil {
		return fmt.Errorf("erro ao marcar convite como utilizado: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}
	if affected == 0 {
		return ErrJaUtilizado
	}

	return nil
}
