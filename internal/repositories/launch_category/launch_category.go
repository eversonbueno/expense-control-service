package launch_category

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type LaunchCategory interface {
	ListLaunchCategorys(ctx context.Context) ([]*entity.LaunchCategory, error)
}

type launchCategory struct {
	db *sql.DB
}

func New(db *sql.DB) LaunchCategory  {
	return &launchCategory{db: db}
}

func (t launchCategory) ListLaunchCategorys(ctx context.Context) ([]*entity.LaunchCategory, error)  {
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
