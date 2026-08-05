package payment_methods

import (
	"net/http"

	paymentMethodsService "expense-control-service/internal/services/payment_methods"

	"github.com/gin-gonic/gin"
)

type PaymentMethods interface {
	ListarFormasPagamento(c *gin.Context)
}

type paymentMethods struct {
	service paymentMethodsService.PaymentMethods
}

func New(service paymentMethodsService.PaymentMethods) PaymentMethods {
	return &paymentMethods{service: service}
}

func (h *paymentMethods) ListarFormasPagamento(c *gin.Context) {
	formas, err := h.service.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao listar formas de pagamento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": formas, "error": nil})
}
