package contas

import (
	"errors"
	"net/http"
	"strconv"

	"expense-control-service/internal/entity"
	"expense-control-service/internal/http/middleware"
	"expense-control-service/internal/http/requests"
	contasService "expense-control-service/internal/services/contas"

	"github.com/gin-gonic/gin"
)

type Contas interface {
	Criar(c *gin.Context)
	Listar(c *gin.Context)
	Atualizar(c *gin.Context)
	Desativar(c *gin.Context)
	Reativar(c *gin.Context)
}

type contas struct {
	service contasService.Contas
}

func New(service contasService.Contas) Contas {
	return &contas{service: service}
}

func grupoFamiliarID(c *gin.Context) (int, bool) {
	raw, ok := c.Get(middleware.ContextKeyGrupoFamiliarID)
	if !ok {
		return 0, false
	}
	id, ok := raw.(int)
	return id, ok
}

func (h *contas) Criar(c *gin.Context) {
	grupoID, ok := grupoFamiliarID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	var req requests.CriarContaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "requisição inválida"})
		return
	}

	conta, err := h.service.Criar(c.Request.Context(), grupoID, req.Nome, req.Tipo, req.FechamentoCartao, req.SaldoInicial)
	if err != nil {
		mapErro(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": conta, "error": nil})
}

func (h *contas) Listar(c *gin.Context) {
	grupoID, ok := grupoFamiliarID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	lista, err := h.service.Listar(c.Request.Context(), grupoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao listar contas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lista, "error": nil})
}

func (h *contas) Atualizar(c *gin.Context) {
	grupoID, ok := grupoFamiliarID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "id inválido"})
		return
	}

	var req requests.AtualizarContaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "requisição inválida"})
		return
	}

	conta, err := h.service.Atualizar(c.Request.Context(), grupoID, id, req.Nome, req.Tipo, req.FechamentoCartao, req.SaldoInicial)
	if err != nil {
		mapErro(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conta, "error": nil})
}

func (h *contas) Desativar(c *gin.Context) {
	h.alterarAtivo(c, false)
}

func (h *contas) Reativar(c *gin.Context) {
	h.alterarAtivo(c, true)
}

func (h *contas) alterarAtivo(c *gin.Context, ativo bool) {
	grupoID, ok := grupoFamiliarID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "id inválido"})
		return
	}

	var conta *entity.Conta
	if ativo {
		conta, err = h.service.Reativar(c.Request.Context(), grupoID, id)
	} else {
		conta, err = h.service.Desativar(c.Request.Context(), grupoID, id)
	}
	if err != nil {
		mapErro(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conta, "error": nil})
}

func mapErro(c *gin.Context, err error) {
	switch {
	case errors.Is(err, contasService.ErrTipoInvalido), errors.Is(err, contasService.ErrFechamentoCartaoInvalido):
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": err.Error()})
	case errors.Is(err, contasService.ErrContaNaoEncontrada):
		c.JSON(http.StatusNotFound, gin.H{"data": nil, "error": err.Error()})
	case errors.Is(err, contasService.ErrAcessoNegado):
		c.JSON(http.StatusForbidden, gin.H{"data": nil, "error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao processar conta"})
	}
}
