package payment_methods

import (
	"context"
	"database/sql"
	"expense-control-service/internal/entity"
	"fmt"
)

type PaymentMethods interface {
	ListPaymentMethodById(ctx context.Context, id int) (*entity.PaymentMethods, error)
	Exists(ctx context.Context, id int) (bool, error)
}

type paymentMethods struct {
	db *sql.DB
}

func New(db *sql.DB) PaymentMethods {
	return &paymentMethods{db: db}
}

func (t paymentMethods) ListPaymentMethodById(ctx context.Context, id int) (*entity.PaymentMethods, error) {
	var paymentMethod entity.PaymentMethods
	row := t.db.QueryRowContext(ctx, ListPaymentMethods, id)
	err := row.Scan(
		&paymentMethod.ID,
		&paymentMethod.Descricao,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao scanear payment methods: %v", err)
	}

	return &paymentMethod, nil
}

func (t paymentMethods) Exists(ctx context.Context, id int) (bool, error) {
	var count int
	row := t.db.QueryRowContext(ctx, CountPaymentMethodById, id)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("erro ao verificar forma de pagamento: %v", err)
	}
	return count > 0, nil
}
