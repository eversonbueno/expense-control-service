package entity

type Launches struct {
	ID                  uint    `json:"id"`
	Usuario             int     `json:"usuario"`
	FormaPagamento      int     `json:"forma_pagamento"`
	TipoLancamento      int     `json:"tipo_lancamento"`
	CategoriaLancamento int     `json:"categoria_lancamento"`
	Mes                 int     `json:"mes"`
	Ano                 int     `json:"ano"`
	Parcelado           int     `json:"parcelado"`
	ParceladoQuantidade int     `json:"parcelado_quantidade"`
	Descricao           string  `json:"descricao"`
	Valor               float64 `json:"valor"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}
