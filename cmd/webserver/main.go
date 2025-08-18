package main

import (
	"context"
	"expense-control-service/internal/container"
)

func main()  {
	ctx := context.Background()

	newContainer := container.New(ctx)

	sr := router
}
