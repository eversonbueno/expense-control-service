package launches

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"expense-control-service/internal/http/middleware"
	"expense-control-service/internal/http/requests"
	launchesService "expense-control-service/internal/services/launches"

	"github.com/gin-gonic/gin"
)

const dataLayout = "2006-01-02"

type Launches interface {
	ListLaunches(c *gin.Context)
	CreateLauches(c *gin.Context)
	UpdateLaunches(c *gin.Context)
	DeleteLaunches(c *gin.Context)
}

type launches struct {
	launchesService launchesService.Launches
}

func New(launchesService launchesService.Launches) Launches {
	return &launches{launchesService: launchesService}
}

func grupoFamiliarIDFromContext(c *gin.Context) (int, bool) {
	raw, ok := c.Get(middleware.ContextKeyGrupoFamiliarID)
	if !ok {
		return 0, false
	}
	id, ok := raw.(int)
	return id, ok
}

func usuarioIDFromContext(c *gin.Context) (int, bool) {
	raw, ok := c.Get(middleware.ContextKeyUsuarioID)
	if !ok {
		return 0, false
	}
	id, ok := raw.(int)
	return id, ok
}

func parseOptionalIntQuery(c *gin.Context, key string) (*int, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// ListLaunches implementa RN-12 (isolamento por grupo familiar) e RN-13 (filtro opcional por mês/ano).
func (l *launches) ListLaunches(c *gin.Context) {
	grupoID, ok := grupoFamiliarIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	mes, err := parseOptionalIntQuery(c, "mes")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "mes inválido"})
		return
	}
	ano, err := parseOptionalIntQuery(c, "ano")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "ano inválido"})
		return
	}

	lancamentos, err := l.launchesService.ListLaunches(c.Request.Context(), grupoID, mes, ano)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao listar lançamentos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lancamentos, "error": nil})
}

func (l *launches) CreateLauches(c *gin.Context) {
	grupoID, ok := grupoFamiliarIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}
	usuarioLogadoID, ok := usuarioIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	var req requests.CriarLancamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "requisição inválida"})
		return
	}

	data, err := time.Parse(dataLayout, req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "data inválida, use o formato AAAA-MM-DD"})
		return
	}

	lancamento, err := l.launchesService.Criar(
		c.Request.Context(), grupoID, usuarioLogadoID, req.ContaID, req.Tipo,
		req.CategoriaID, req.FormaPagamentoID, req.Descricao, req.Valor, data,
	)
	if err != nil {
		mapErro(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": lancamento, "error": nil})
}

func (l *launches) UpdateLaunches(c *gin.Context) {
	grupoID, ok := grupoFamiliarIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "id inválido"})
		return
	}

	var req requests.AtualizarLancamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "requisição inválida"})
		return
	}

	data, err := time.Parse(dataLayout, req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "data inválida, use o formato AAAA-MM-DD"})
		return
	}

	lancamento, err := l.launchesService.Atualizar(
		c.Request.Context(), grupoID, id, req.ContaID, req.Tipo,
		req.CategoriaID, req.FormaPagamentoID, req.Descricao, req.Valor, data,
	)
	if err != nil {
		mapErro(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": lancamento, "error": nil})
}

func (l *launches) DeleteLaunches(c *gin.Context) {
	grupoID, ok := grupoFamiliarIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "id inválido"})
		return
	}

	if err := l.launchesService.Excluir(c.Request.Context(), grupoID, id); err != nil {
		mapErro(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "lançamento excluído", "error": nil})
}

func mapErro(c *gin.Context, err error) {
	switch {
	case errors.Is(err, launchesService.ErrTipoInvalido),
		errors.Is(err, launchesService.ErrValorInvalido),
		errors.Is(err, launchesService.ErrContaInvalida),
		errors.Is(err, launchesService.ErrCategoriaInvalida),
		errors.Is(err, launchesService.ErrFormaPagamentoInvalida):
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": err.Error()})
	case errors.Is(err, launchesService.ErrLancamentoNaoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{"data": nil, "error": err.Error()})
	case errors.Is(err, launchesService.ErrAcessoNegado):
		c.JSON(http.StatusForbidden, gin.H{"data": nil, "error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao processar lançamento"})
	}
}
