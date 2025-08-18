package entity

type PaymentMethods struct {
	ID        uint   `json:"id"`
	Descricao string `json:"descricao"`
}

var (
	PaymentMethodsDebito = PaymentMethods{
		ID:        1,
		Descricao: "DEBITO",
	}
	PaymentMethodsCredito = PaymentMethods{
		ID:        2,
		Descricao: "CREDITO",
	}
	PaymentMethodsPix = PaymentMethods{
		ID:        3,
		Descricao: "PIX",
	}
	PaymentMethodsDinheiro = PaymentMethods{
		ID:        4,
		Descricao: "DINHEIRO",
	}
	PaymentMethodsOutros = PaymentMethods{
		ID:        5,
		Descricao: "OUTROS",
	}
)
