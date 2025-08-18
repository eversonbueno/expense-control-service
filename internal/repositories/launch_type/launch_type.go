package launch_type

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type LauchType interface {
	ListLauchTypes(ctx context.Context) ([]*entity.LauchType, error)
}

type lauchType struct {
	db *sql.DB
}

func New(db *sql.DB) LauchType  {
	return &lauchType{db: db}
}

func (t lauchType) ListLauchTypes(ctx context.Context) ([]*entity.LauchType, error)  {
	rows, err := t.db.QueryContext(ctx, ListLaunchType)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %v", err)
	}

	defer rows.Close()

	var lauchTypes []*entity.LauchType
	for rows.Next() {
		var lauchType entity.LauchType
		err := rows.Scan(
			&lauchType.ID,
			&lauchType.Descricao,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear lauch types: %v", err)
		}

		lauchTypes = append(lauchTypes, &lauchType)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar linhas: %v", err)
	}

	return lauchTypes, nil
}
