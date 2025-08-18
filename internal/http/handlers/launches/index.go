package launches

import (
	launchesService "expense-control-service/internal/services/launches"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Launches interface {
	ListLaunches(c *gin.Context)
}

type launches struct {
	launchesService launchesService.Launches
}

func New(launchesService launchesService.Launches) Launches {
	return &launches{launchesService: launchesService}
}

type ReponseLaunches struct {
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
}

func (l launches) ListLaunches(c *gin.Context)  {
	lauches, err := l.launchesService.ListLaunches(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao listar lançamentos"})
		return
	}

	resp := make([]ReponseLaunches, len(lauches))
	for i, r := range lauches {
		resp[i] = ReponseLaunches{
			Usuario:             r.Usuario,
			FormaPagamento:      r.FormaPagamento,
			TipoLancamento:      r.TipoLancamento,
			CategoriaLancamento: r.CategoriaLancamento,
			Mes:                 r.Mes,
			Ano:                 r.Ano,
			Parcelado:           r.Parcelado,
			ParceladoQuantidade: r.ParceladoQuantidade,
			Descricao:           r.Descricao,
			Valor:               25.25,
		}
	}

	c.JSON(http.StatusOK, resp)
}


