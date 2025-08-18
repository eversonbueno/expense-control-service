package payment_methods

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type PaymentMethods interface {
	ListLauchTypes(ctx context.Context) ([]*entity.PaymentMethods, error)
}

type paymentMethods struct {
	db *sql.DB
}

func New(db *sql.DB) PaymentMethods  {
	return &paymentMethods{db: db}
}

func (t paymentMethods) ListLauchTypes(ctx context.Context) ([]*entity.PaymentMethods, error)  {
	rows, err := t.db.QueryContext(ctx, ListPaymentMethods)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %v", err)
	}

	defer rows.Close()

	var paymentMethods []*entity.PaymentMethods
	for rows.Next() {
		var paymentMethod entity.PaymentMethods
		err := rows.Scan(
			&paymentMethod.ID,
			&paymentMethod.Descricao,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear payment methods: %v", err)
		}

		paymentMethods = append(paymentMethods, &paymentMethod)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar linhas: %v", err)
	}

	return paymentMethods, nil
}
