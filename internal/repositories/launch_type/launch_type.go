package launch_type

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type LauchType interface {
	ListLauchTypesById(ctx context.Context, id int) (*entity.LauchType, error)
}

type lauchType struct {
	db *sql.DB
}

func New(db *sql.DB) LauchType  {
	return &lauchType{db: db}
}

func (t lauchType) ListLauchTypesById(ctx context.Context, id int) (*entity.LauchType, error)  {
	var lauchType entity.LauchType
	row := t.db.QueryRowContext(ctx, ListLaunchType, id)
	err := row.Scan(
		&lauchType.ID,
		&lauchType.Descricao,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao scanear lauch type: %v", err)
	}

	return &lauchType, nil
}
