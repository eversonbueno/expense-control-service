package router

import (
	"expense-control-service/internal/container"
	"expense-control-service/internal/http/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupRouter(container *container.Container) *gin.Engine  {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", handlers.RedirectToAlive)
	r.GET("/health-check/alive", handlers.Alive)
	r.GET("/health-check/status", handlers.Status)

	v1 := r.Group("/api/v1")
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
