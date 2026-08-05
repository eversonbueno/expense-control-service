package launches

import (
	"context"
	"database/sql"
	"errors"
	"expense-control-service/internal/entity"
	"fmt"
)

type Launches interface {
	ListLaunches(ctx context.Context, grupoFamiliarID int, mes, ano *int) ([]*entity.Launches, error)
	FindByID(ctx context.Context, id int) (*entity.Launches, error)
	Create(ctx context.Context, launch *entity.Launches) (*entity.Launches, error)
	Update(ctx context.Context, launch *entity.Launches) (*entity.Launches, error)
	SoftDelete(ctx context.Context, id int) error
}

type launches struct {
	db *sql.DB
}

func New(db *sql.DB) Launches {
	return &launches{db: db}
}

func scanLaunch(row *sql.Row, l *entity.Launches) error {
	return row.Scan(
		&l.ID,
		&l.Usuario,
		&l.GrupoFamiliarID,
		&l.ContaID,
		&l.FormaPagamento,
		&l.TipoLancamento,
		&l.CategoriaLancamento,
		&l.Mes,
		&l.Ano,
		&l.Data,
		&l.Parcelado,
		&l.ParceladoQuantidade,
		&l.Descricao,
		&l.Valor,
		&l.Excluido,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
}

// ListLaunches aplica RN-12 (isolamento por grupo familiar, exclui excluídos) e RN-13
// (filtro opcional por mês/ano — mes/ano nil desativam o respectivo filtro).
func (t launches) ListLaunches(ctx context.Context, grupoFamiliarID int, mes, ano *int) ([]*entity.Launches, error) {
	rows, err := t.db.QueryContext(ctx, ListLaunches, grupoFamiliarID, mes, mes, ano, ano)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %v", err)
	}

	defer rows.Close()

	var launches []*entity.Launches
	for rows.Next() {
		var launch entity.Launches
		err := rows.Scan(
			&launch.ID,
			&launch.Usuario,
			&launch.GrupoFamiliarID,
			&launch.ContaID,
			&launch.FormaPagamento,
			&launch.TipoLancamento,
			&launch.CategoriaLancamento,
			&launch.Mes,
			&launch.Ano,
			&launch.Data,
			&launch.Parcelado,
			&launch.ParceladoQuantidade,
			&launch.Descricao,
			&launch.Valor,
			&launch.Excluido,
			&launch.CreatedAt,
			&launch.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear lançamento: %v", err)
		}

		launches = append(launches, &launch)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar linhas: %v", err)
	}

	return launches, nil
}

// FindByID retorna nil, nil quando nenhum lançamento é encontrado com o id informado.
func (t launches) FindByID(ctx context.Context, id int) (*entity.Launches, error) {
	var launch entity.Launches
	row := t.db.QueryRowContext(ctx, FindLaunchById, id)

	if err := scanLaunch(row, &launch); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar lançamento: %v", err)
	}

	return &launch, nil
}

func (t launches) Create(ctx context.Context, launch *entity.Launches) (*entity.Launches, error) {
	result, err := t.db.ExecContext(ctx, CreateLaunch,
		launch.Usuario, launch.GrupoFamiliarID, launch.ContaID, launch.FormaPagamento,
		launch.TipoLancamento, launch.CategoriaLancamento, launch.Data, launch.Descricao, launch.Valor,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar lançamento: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter id do lançamento criado: %v", err)
	}

	return t.FindByID(ctx, int(id))
}

func (t launches) Update(ctx context.Context, launch *entity.Launches) (*entity.Launches, error) {
	_, err := t.db.ExecContext(ctx, UpdateLaunch,
		launch.ContaID, launch.FormaPagamento, launch.TipoLancamento, launch.CategoriaLancamento,
		launch.Data, launch.Descricao, launch.Valor, launch.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar lançamento: %v", err)
	}

	return t.FindByID(ctx, int(launch.ID))
}

func (t launches) SoftDelete(ctx context.Context, id int) error {
	_, err := t.db.ExecContext(ctx, SoftDeleteLaunch, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir lançamento: %v", err)
	}
	return nil
}
