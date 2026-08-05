package entity

import "time"

type Launches struct {
	ID                  uint      `json:"id"`
	Usuario             int       `json:"usuario"`
	GrupoFamiliarID     int       `json:"grupo_familiar_id"`
	ContaID             int       `json:"conta_id"`
	FormaPagamento      int       `json:"forma_pagamento"`
	TipoLancamento      int       `json:"tipo_lancamento"`
	CategoriaLancamento int       `json:"categoria_lancamento"`
	Mes                 *int      `json:"mes"`
	Ano                 *int      `json:"ano"`
	Data                time.Time `json:"data"`
	Parcelado           int       `json:"parcelado"`
	ParceladoQuantidade int       `json:"parcelado_quantidade"`
	Descricao           string    `json:"descricao"`
	Valor               float64   `json:"valor"`
	Excluido            bool      `json:"excluido"`
	CreatedAt           string    `json:"created_at"`
	UpdatedAt           string    `json:"updated_at"`
}
