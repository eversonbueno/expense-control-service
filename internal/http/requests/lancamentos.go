package requests

type CriarLancamentoRequest struct {
	ContaID          int     `json:"conta_id" binding:"required"`
	Tipo             string  `json:"tipo" binding:"required"`
	CategoriaID      int     `json:"categoria_id" binding:"required"`
	FormaPagamentoID int     `json:"forma_pagamento_id" binding:"required"`
	Descricao        string  `json:"descricao" binding:"required"`
	Valor            float64 `json:"valor"`
	Data             string  `json:"data" binding:"required"`
}

type AtualizarLancamentoRequest struct {
	ContaID          int     `json:"conta_id" binding:"required"`
	Tipo             string  `json:"tipo" binding:"required"`
	CategoriaID      int     `json:"categoria_id" binding:"required"`
	FormaPagamentoID int     `json:"forma_pagamento_id" binding:"required"`
	Descricao        string  `json:"descricao" binding:"required"`
	Valor            float64 `json:"valor"`
	Data             string  `json:"data" binding:"required"`
}
