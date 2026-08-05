package payment_methods

import (
	"context"
	"expense-control-service/internal/entity"
	paymentMethodsRepo "expense-control-service/internal/repositories/payment_methods"
)

type PaymentMethods interface {
	ListPaymentMethodsById(ctx context.Context, id int) (*entity.PaymentMethods, error)
	ListAll(ctx context.Context) ([]*entity.PaymentMethods, error)
}

type paymentMethods struct {
	paymentMethodRepo paymentMethodsRepo.PaymentMethods
}

func New(
	paymentMethodRepo paymentMethodsRepo.PaymentMethods,
) PaymentMethods {
	return &paymentMethods{paymentMethodRepo: paymentMethodRepo}
}

func (u *paymentMethods) ListPaymentMethodsById(ctx context.Context, id int) (*entity.PaymentMethods, error) {
	paymentMethod, err := u.paymentMethodRepo.ListPaymentMethodById(ctx, id)
	if err != nil {
		return nil, err
	}

	return paymentMethod, nil
}

func (u *paymentMethods) ListAll(ctx context.Context) ([]*entity.PaymentMethods, error) {
	paymentMethods, err := u.paymentMethodRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return paymentMethods, nil
}
