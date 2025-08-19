package main

import (
	"context"
	"expense-control-service/internal/container"
	"expense-control-service/internal/http/router"
)

func main()  {
	ctx := context.Background()

	newContainer := container.New(ctx)

	sr := router.SetupRouter(newContainer)
	router.RunServer(sr, ":9000")
}
