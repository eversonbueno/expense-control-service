package payment_methods

var (
	ListPaymentMethods = `
		SELECT
			id,
			descricao
		FROM formas_pagamento
		WHERE id = ?
	`

	CountPaymentMethodById = `
		SELECT COUNT(*)
		FROM formas_pagamento
		WHERE id = ?
	`
)
