package launches

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type Launches interface {
	ListLaunches(ctx context.Context) ([]*entity.Launches, error)
}

type launches struct {
	db *sql.DB
}

func New(db *sql.DB) Launches  {
	return &launches{db: db}
}

func (t launches) ListLaunches(ctx context.Context) ([]*entity.Launches, error)  {
	rows, err := t.db.QueryContext(ctx, ListLaunches)
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
			&launch.FormaPagamento,
			&launch.TipoLancamento,
			&launch.CategoriaLancamento,
			&launch.Mes,
			&launch.Ano,
			&launch.Parcelado,
			&launch.ParceladoQuantidade,
			&launch.Descricao,
			&launch.Valor,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear lauch types: %v", err)
		}

		launches = append(launches, &launch)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar linhas: %v", err)
	}

	return launches, nil
}

