package main

import (
	"expense-control-service/internal/commands"
	"expense-control-service/pkg/console"
	"fmt"
	"os"
)

func init()  {
	console.AddCommand(
		"help",
		"Lista todos os comandos disponiveis",
		console.Help,
	)

	console.AddCommand(
		"exec:list-users",
		"Lista os usuários cadastrados",
		commands.UsersListCommand,
	)
}

func main() {
	if len(os.Args) < 2 {
		console.PrintError("Uso: go run ./cmd/cli/main.go help")
		return
	}

	name, args := os.Args[1], os.Args[2:]
	cmd, ok := console.GetCommand(name)
	if !ok {
		fmt.Printf("Comando não encontrado: %s\n", name)
		return
	}

	opts := console.ParseOptions(args)
	for k, v := range opts {
		cmd.Flags().String(k, v, "")
	}

	cmd.SetArgs(args)
	cmd.Execute()
}
