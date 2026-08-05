package router

import (
	"expense-control-service/internal/container"
	"expense-control-service/internal/http/handlers"
	"expense-control-service/internal/http/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupRouter(container *container.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware(container.Config.CORSAllowedOrigins))

	r.GET("/", handlers.RedirectToAlive)
	r.GET("/health-check/alive", handlers.Alive)
	r.GET("/health-check/status", handlers.Status)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/registrar", container.Handler.Auth.Registrar)
		authGroup.POST("/login", container.Handler.Auth.Login)
		authGroup.POST("/convite", middleware.AuthMiddleware(container.Config.JWTSecret), container.Handler.Auth.GerarConvite)
	}

	v1 := r.Group("/api/v1", middleware.AuthMiddleware(container.Config.JWTSecret))
	{
		lancamentos := v1.Group("/lancamentos")
		{
			lancamentos.GET("", container.Handler.ListLaunches)
			lancamentos.POST("", container.Handler.CreateLauches)
			lancamentos.PUT("/:id", container.Handler.UpdateLaunches)
			lancamentos.DELETE("/:id", container.Handler.DeleteLaunches)
		}

		contas := v1.Group("/contas")
		{
			contas.POST("", container.Handler.Contas.Criar)
			contas.GET("", container.Handler.Contas.Listar)
			contas.PUT("/:id", container.Handler.Contas.Atualizar)
			contas.PATCH("/:id/desativar", container.Handler.Contas.Desativar)
			contas.PATCH("/:id/reativar", container.Handler.Contas.Reativar)
		}

		v1.GET("/categorias", container.Handler.Categorias.ListarCategorias)
		v1.GET("/formas-pagamento", container.Handler.FormasPagamento.ListarFormasPagamento)
	}

	return r
}

func RunServer(router *gin.Engine, addr string) {
	fmt.Printf("Servidor rodando em http://127.0.0.1%s\n", addr)
	_ = router.Run(addr)
}
