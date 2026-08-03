package commands

import (
	"context"
	"expense-control-service/internal/container"
	"expense-control-service/pkg/console"
	"fmt"
	"github.com/spf13/cobra"
)

func UsersListCommand(cmd *cobra.Command, args []string){
	console.PrintSuccess("\nIniciando o comando 'exec:purchase-order-bip-received'")

	ctx := context.Background()
	newContainer := container.New(ctx)

	users, err := newContainer.Service.User.ListUsers(ctx)
	if err != nil {
		_ = fmt.Errorf("erro ao buscar Usuários", nil, err)
	}

	for _, user := range users {
		fmt.Printf("ID: %d, Nome: %s\n", user.ID, user.Nome)
	}
}
