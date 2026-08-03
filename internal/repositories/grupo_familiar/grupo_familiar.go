package grupo_familiar

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type GrupoFamiliar interface {
	Create(ctx context.Context) (*entity.GrupoFamiliar, error)
}

type grupoFamiliar struct {
	db *sql.DB
}

func New(db *sql.DB) GrupoFamiliar {
	return &grupoFamiliar{db: db}
}

func (g *grupoFamiliar) Create(ctx context.Context) (*entity.GrupoFamiliar, error) {
	result, err := g.db.ExecContext(ctx, CreateGrupoFamiliar)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar grupo familiar: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter id do grupo familiar criado: %v", err)
	}

	var grupo entity.GrupoFamiliar
	row := g.db.QueryRowContext(ctx, GetGrupoFamiliarById, id)
	if err := row.Scan(&grupo.ID, &grupo.CriadoEm); err != nil {
		return nil, fmt.Errorf("erro ao buscar grupo familiar criado: %v", err)
	}

	return &grupo, nil
}
