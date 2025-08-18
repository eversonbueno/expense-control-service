package container

import (
	"context"
	"database/sql"
	"expense-control-service/configs"
	"expense-control-service/internal/http/handlers/launches"
	launchesHandler "expense-control-service/internal/http/handlers/launches"
	launchCategoryRepo "expense-control-service/internal/repositories/launch_category"
	launchTypeRepo "expense-control-service/internal/repositories/launch_type"
	launchesRepo "expense-control-service/internal/repositories/launches"
	paymentMethodsRepo "expense-control-service/internal/repositories/payment_methods"
	usersRepo "expense-control-service/internal/repositories/users"
	launchesService "expense-control-service/internal/services/launches"
	userServ "expense-control-service/internal/services/user"
	"fmt"
	"time"
)

type Repositories struct {
	Users          usersRepo.Users
	LaunchType     launchTypeRepo.LauchType
	PaymentMethods paymentMethodsRepo.PaymentMethods
	LaunchCategory launchCategoryRepo.LaunchCategory
	Launches       launchesRepo.Launches
}
type Services struct {
	User userServ.User
	Launches launchesService.Launches
}
type Handlers struct{
	launches.Launches
}
type Container struct {
	Config     configs.AppConfig
	DB         *sql.DB
	Repository Repositories
	Service    Services
	Handler    Handlers
}

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
	newLaunchTypeRepo := launchTypeRepo.New(db)
	newPaymentMethodsRepo := paymentMethodsRepo.New(db)
	newLaunchCategoryRepo := launchCategoryRepo.New(db)
	newLaunchesRepo := launchesRepo.New(db)

	//Services
	newUserServ := userServ.New(newUsersRepo)
	newLaunchesService := launchesService.New(newLaunchesRepo)

	//Handlers
	newLaunchesHandler := launchesHandler.New(newLaunchesService)

	return &Container{
		Config: *cfg,
		DB:     db,
		Repository: Repositories{
			Users:          newUsersRepo,
			LaunchType:     newLaunchTypeRepo,
			PaymentMethods: newPaymentMethodsRepo,
			LaunchCategory: newLaunchCategoryRepo,
			Launches:       newLaunchesRepo,
		},
		Service: Services{
			User: newUserServ,
			Launches: newLaunchesService,
		},
		Handler: Handlers{
			newLaunchesHandler,
		},
	}

}
