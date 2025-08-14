package container

import (
	"context"
	"database/sql"
	"expense-control-service/configs"
	usersRepo "expense-control-service/internal/repositories/users"
	userServ "expense-control-service/internal/services/user"
	"fmt"
	"time"
)

type Repositories struct {
	Users usersRepo.Users
}
type Services struct {
	User userServ.User
}
type Handlers struct {}
type Container struct {
	Config            configs.AppConfig
	DB                *sql.DB
	Repository        Repositories
	Service           Services
}
type Consumers struct {}

var New = newContainer
func newContainer(ctx context.Context) *Container {
	cfg := new(configs.AppConfig)
	cfg.Bootstrap(ctx)
	start := time.Now()
	cfg.MySqlLatency = time.Since(start).Milliseconds()

	db, err := cfg.Mysql.Connect()
	if err != nil {
		fmt.Println("Erro MySQL:", err)
	}

	//Repositories
	newUsersRepo := usersRepo.New(db)

	//Services
	newUserServ := userServ.New(newUsersRepo)

	return &Container{
		Config: *cfg,
		DB:     db,
		Repository: Repositories{
			Users: newUsersRepo,
		},
		Service: Services{
			User: newUserServ,
		},
	}


}