package launch_category

import (
	"context"
	"database/sql"
	"errors"
	"expense-control-service/internal/entity"
	"fmt"
)

type LaunchCategory interface {
	ListLaunchCategorys(ctx context.Context) ([]*entity.LaunchCategory, error)
	FindByID(ctx context.Context, id int) (*entity.LaunchCategory, error)
}

type launchCategory struct {
	db *sql.DB
}

func New(db *sql.DB) LaunchCategory {
	return &launchCategory{db: db}
}

func (t launchCategory) ListLaunchCategorys(ctx context.Context) ([]*entity.LaunchCategory, error) {
	rows, err := t.db.QueryContext(ctx, ListLaunchCategory)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %v", err)
	}

	defer rows.Close()

	var lauchCategorys []*entity.LaunchCategory
	for rows.Next() {
		var lauchCategory entity.LaunchCategory
		err := rows.Scan(
			&lauchCategory.ID,
			&lauchCategory.Descricao,
			&lauchCategory.TipoLancamentoID,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear lauch types: %v", err)
		}

		lauchCategorys = append(lauchCategorys, &lauchCategory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar linhas: %v", err)
	}

	return lauchCategorys, nil
}

// FindByID retorna nil, nil quando nenhuma categoria é encontrada com o id informado.
func (t launchCategory) FindByID(ctx context.Context, id int) (*entity.LaunchCategory, error) {
	var category entity.LaunchCategory
	row := t.db.QueryRowContext(ctx, FindLaunchCategoryById, id)

	if err := row.Scan(&category.ID, &category.Descricao, &category.TipoLancamentoID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar categoria de lançamento: %v", err)
	}

	return &category, nil
}
