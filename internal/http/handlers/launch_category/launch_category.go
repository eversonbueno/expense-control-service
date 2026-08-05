package launch_category

import (
	"net/http"

	launchCategoryService "expense-control-service/internal/services/launch_category"

	"github.com/gin-gonic/gin"
)

type LaunchCategory interface {
	ListarCategorias(c *gin.Context)
}

type launchCategory struct {
	service launchCategoryService.LaunchCategory
}

func New(service launchCategoryService.LaunchCategory) LaunchCategory {
	return &launchCategory{service: service}
}

func (h *launchCategory) ListarCategorias(c *gin.Context) {
	categorias, err := h.service.ListLaunchCategorys(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao listar categorias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categorias, "error": nil})
}
