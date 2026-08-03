package auth

import (
	"errors"
	"expense-control-service/internal/http/middleware"
	"expense-control-service/internal/http/requests"
	authService "expense-control-service/internal/services/auth"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Auth interface {
	Registrar(c *gin.Context)
	Login(c *gin.Context)
	GerarConvite(c *gin.Context)
}

type auth struct {
	service authService.Auth
}

func New(service authService.Auth) Auth {
	return &auth{service: service}
}

type registrarResponse struct {
	ID              int    `json:"id"`
	Nome            string `json:"nome"`
	Email           string `json:"email"`
	GrupoFamiliarID int    `json:"grupo_familiar_id"`
	CriadoEm        string `json:"criado_em"`
}

type loginResponse struct {
	Token    string `json:"token"`
	ExpiraEm string `json:"expira_em"`
}

type conviteResponse struct {
	Codigo   string `json:"codigo"`
	ExpiraEm string `json:"expira_em"`
}

func (h *auth) Registrar(c *gin.Context) {
	var req requests.RegistrarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "requisição inválida"})
		return
	}

	user, err := h.service.Registrar(c.Request.Context(), req.Nome, req.Email, req.Senha, req.CodigoConvite)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrEmailJaCadastrado):
			c.JSON(http.StatusConflict, gin.H{"data": nil, "error": err.Error()})
		case errors.Is(err, authService.ErrSenhaCurta), errors.Is(err, authService.ErrConvite):
			c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao registrar usuário"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": registrarResponse{
			ID:              user.ID,
			Nome:            user.Nome,
			Email:           user.Email,
			GrupoFamiliarID: user.GrupoFamiliarID,
			CriadoEm:        user.CreatedAt.Format(time.RFC3339),
		},
		"error": nil,
	})
}

func (h *auth) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "requisição inválida"})
		return
	}

	token, expiraEm, err := h.service.Login(c.Request.Context(), req.Email, req.Senha)
	if err != nil {
		if errors.Is(err, authService.ErrCredenciaisInvalidas) {
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": err.Error()})
			return
		}
		log.Printf("erro interno ao autenticar %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao autenticar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": loginResponse{
			Token:    token,
			ExpiraEm: expiraEm.Format(time.RFC3339),
		},
		"error": nil,
	})
}

func (h *auth) GerarConvite(c *gin.Context) {
	raw, ok := c.Get(middleware.ContextKeyGrupoFamiliarID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}
	grupoFamiliarID, ok := raw.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "não autenticado"})
		return
	}

	convite, err := h.service.GerarConvite(c.Request.Context(), grupoFamiliarID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "erro ao gerar convite"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": conviteResponse{
			Codigo:   convite.Codigo,
			ExpiraEm: convite.ExpiraEm.Format(time.RFC3339),
		},
		"error": nil,
	})
}
