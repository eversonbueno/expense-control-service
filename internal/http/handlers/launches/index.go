package launches

import (
	lauchesTypeService "expense-control-service/internal/services/lauches_type"
	launchesService "expense-control-service/internal/services/launches"
	paymentMethodsService "expense-control-service/internal/services/payment_methods"
	userService "expense-control-service/internal/services/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Launches interface {
	ListLaunches(c *gin.Context)
	CreateLauches(c *gin.Context)
}

type launches struct {
	launchesService   launchesService.Launches
	usersService      userService.User
	paymentService    paymentMethodsService.PaymentMethods
	laucheTypeService lauchesTypeService.LauchType
}

func New(
	launchesService launchesService.Launches,
	usersService userService.User,
	paymentService paymentMethodsService.PaymentMethods,
	laucheTypeService lauchesTypeService.LauchType,
) Launches {
	return &launches{
		launchesService:   launchesService,
		usersService:      usersService,
		paymentService:    paymentService,
		laucheTypeService: laucheTypeService,
	}
}

type ReponseLaunches struct {
	Usuario             string  `json:"usuario"`
	FormaPagamento      string  `json:"forma_pagamento"`
	TipoLancamento      string  `json:"tipo_lancamento"`
	CategoriaLancamento int     `json:"categoria_lancamento"`
	Mes                 string     `json:"mes"`
	Ano                 int     `json:"ano"`
	Parcelado           int     `json:"parcelado"`
	ParceladoQuantidade int     `json:"parcelado_quantidade"`
	Descricao           string  `json:"descricao"`
	Valor               float64 `json:"valor"`
}

func (l launches) ListLaunches(c *gin.Context) {
	mesesMap := map[int]string{
		1:  "Janeiro",
		2:  "Fevereiro",
		3:  "Março",
		4:  "Abril",
		5:  "Maio",
		6:  "Junho",
		7:  "Julho",
		8:  "Agosto",
		9:  "Setembro",
		10: "Outubro",
		11: "Novembro",
		12: "Dezembro",
	}

	lauches, err := l.launchesService.ListLaunches(c.Request.Context())
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao listar lançamentos"})
		return
	}

	resp := make([]ReponseLaunches, len(lauches))
	for i, r := range lauches {
		user, err := l.usersService.ListUserById(c.Request.Context(), int(r.ID))
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "falha ao buscar dados do usuário"})
			return
		}

		paymentMethod, err := l.paymentService.ListPaymentMethodsById(c.Request.Context(), int(r.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "falha ao buscar dados do método de pagamento"})
			return
		}

		lauchType, err := l.laucheTypeService.ListLauchTypesById(c.Request.Context(), int(r.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "falha ao buscar dados do tipo de lançamento"})
			return
		}

		nomeMes, ok := mesesMap[r.Mes];
		if !ok {

		}

		resp[i] = ReponseLaunches{
			Usuario:             user.Nome,
			FormaPagamento:      paymentMethod.Descricao,
			TipoLancamento:      lauchType.Descricao,
			CategoriaLancamento: r.CategoriaLancamento,
			Mes:                 strconv.Itoa(r.Mes) + "-" + nomeMes,
			Ano:                 r.Ano,
			Parcelado:           r.Parcelado,
			ParceladoQuantidade: r.ParceladoQuantidade,
			Descricao:           r.Descricao,
			Valor:               r.Valor,
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (l launches) CreateLauches(c *gin.Context) {

}