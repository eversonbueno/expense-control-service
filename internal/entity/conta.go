package entity

import "time"

type Conta struct {
	ID               int       `json:"id"`
	GrupoFamiliarID  int       `json:"grupo_familiar_id"`
	Nome             string    `json:"nome"`
	Tipo             string    `json:"tipo"`
	FechamentoCartao *int      `json:"fechamento_cartao"`
	SaldoInicial     float64   `json:"saldo_inicial"`
	Ativo            bool      `json:"ativo"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
