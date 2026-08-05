package container

import (
	"context"
	"database/sql"
	"expense-control-service/configs"
	authHandler "expense-control-service/internal/http/handlers/auth"
	contasHandler "expense-control-service/internal/http/handlers/contas"
	"expense-control-service/internal/http/handlers/launches"
	launchesHandler "expense-control-service/internal/http/handlers/launches"
	contasRepo "expense-control-service/internal/repositories/contas"
	conviteRepo "expense-control-service/internal/repositories/convite"
	grupoFamiliarRepo "expense-control-service/internal/repositories/grupo_familiar"
	launchCategoryRepo "expense-control-service/internal/repositories/launch_category"
	launchTypeRepo "expense-control-service/internal/repositories/launch_type"
	launchesRepo "expense-control-service/internal/repositories/launches"
	paymentMethodsRepo "expense-control-service/internal/repositories/payment_methods"
	usersRepo "expense-control-service/internal/repositories/users"
	authService "expense-control-service/internal/services/auth"
	contasService "expense-control-service/internal/services/contas"
	lauchTypesService "expense-control-service/internal/services/lauches_type"
	launchesService "expense-control-service/internal/services/launches"
	paymentMethodsServ "expense-control-service/internal/services/payment_methods"
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
	GrupoFamiliar  grupoFamiliarRepo.GrupoFamiliar
	Convite        conviteRepo.Convite
	Contas         contasRepo.Contas
}
type Services struct {
	User userServ.User
	Launches launchesService.Launches
	PaymentMethodsService paymentMethodsServ.PaymentMethods
	LauchTypesService lauchTypesService.LauchType
	Auth authService.Auth
	Contas contasService.Contas
}
type Handlers struct{
	launches.Launches
	Auth authHandler.Auth
	Contas contasHandler.Contas
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
	newGrupoFamiliarRepo := grupoFamiliarRepo.New(db)
	newConviteRepo := conviteRepo.New(db)
	newContasRepo := contasRepo.New(db)

	//Services
	newUserServ := userServ.New(newUsersRepo)
	newLaunchesService := launchesService.New(newLaunchesRepo, newContasRepo, newLaunchCategoryRepo, newPaymentMethodsRepo)
	newPaymentMethodsService := paymentMethodsServ.New(newPaymentMethodsRepo)
	newLauchTypesService := lauchTypesService.New(newLaunchTypeRepo)
	newAuthService := authService.New(newUsersRepo, newGrupoFamiliarRepo, newConviteRepo, cfg.JWTSecret)
	newContasService := contasService.New(newContasRepo)

	//Handlers
	newLaunchesHandler := launchesHandler.New(newLaunchesService)
	newAuthHandler := authHandler.New(newAuthService)
	newContasHandler := contasHandler.New(newContasService)

	return &Container{
		Config: *cfg,
		DB:     db,
		Repository: Repositories{
			Users:          newUsersRepo,
			LaunchType:     newLaunchTypeRepo,
			PaymentMethods: newPaymentMethodsRepo,
			LaunchCategory: newLaunchCategoryRepo,
			Launches:       newLaunchesRepo,
			GrupoFamiliar:  newGrupoFamiliarRepo,
			Convite:        newConviteRepo,
			Contas:         newContasRepo,
		},
		Service: Services{
			User: newUserServ,
			Launches: newLaunchesService,
			PaymentMethodsService: newPaymentMethodsService,
			LauchTypesService: newLauchTypesService,
			Auth: newAuthService,
			Contas: newContasService,
		},
		Handler: Handlers{
			newLaunchesHandler,
			newAuthHandler,
			newContasHandler,
		},
	}

}
