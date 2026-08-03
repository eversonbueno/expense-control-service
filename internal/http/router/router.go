package router

import (
	"expense-control-service/internal/container"
	"expense-control-service/internal/http/handlers"
	"expense-control-service/internal/http/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupRouter(container *container.Container) *gin.Engine  {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

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
		expenseControl := v1.Group("/expense-control")
		{
			expenseControl.GET("", container.Handler.ListLaunches)
			expenseControl.POST("/create", container.Handler.CreateLauches)
		}
	}

	return r
}

func RunServer(router *gin.Engine, addr string) {
	fmt.Printf("Servidor rodando em http://127.0.0.1%s\n", addr)
	_ = router.Run(addr)
}
