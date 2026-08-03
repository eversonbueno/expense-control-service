package configs

import (
	"context"
	"database/sql"
	"expense-control-service/pkg/mysql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	AppName                        string
	AppEnv                         string
	Mysql                          mysql.MySQLInterface
	MySqlLatency                   int64
	JWTSecret                      string
	CORSAllowedOrigins             []string
}

var GlobalConfig AppConfig
var MysqlDB *sql.DB

func (cfg *AppConfig) Bootstrap(ctx context.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		_ = fmt.Errorf("Arquivo .env não encontrado!", nil, err)
		os.Exit(1)
	}

	cfg.AppName = os.Getenv("APP_NAME")
	cfg.AppEnv = os.Getenv("APP_ENV")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.CORSAllowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	cfg.Mysql = mysql.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		5,
		5,
		time.Duration(60),
	)
}
