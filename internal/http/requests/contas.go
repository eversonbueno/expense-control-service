package requests

type CriarContaRequest struct {
	Nome             string   `json:"nome" binding:"required"`
	Tipo             string   `json:"tipo" binding:"required"`
	FechamentoCartao *int     `json:"fechamento_cartao"`
	SaldoInicial     *float64 `json:"saldo_inicial"`
}

type AtualizarContaRequest struct {
	Nome             string  `json:"nome" binding:"required"`
	Tipo             string  `json:"tipo" binding:"required"`
	FechamentoCartao *int    `json:"fechamento_cartao"`
	SaldoInicial     float64 `json:"saldo_inicial"`
}
