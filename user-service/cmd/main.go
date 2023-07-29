package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stas-bukovskiy/wish-scribe/packages/database"
	"github.com/stas-bukovskiy/wish-scribe/packages/httpserver"
	"github.com/stas-bukovskiy/wish-scribe/packages/logger"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/entity"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/handler"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/repository"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/service"
	"os"
	"os/signal"
	"syscall"
)

// @title           User Service API
// @version         1.0
// @description     This is service to authenticate users and verify their JWT tokens
//
// @host localhost:8000
// @BathPath /
//
// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log := logger.New("INFO")
	if err := initConfigs(); err != nil {
		log.Fatal("error occurred while config initializing: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("error occurred while env variables loading: %s", err.Error())
	}

	db, err := database.NewPostgreSQL(database.PostgreSQLConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetBool("db.sslmode"),
		TimeZone: viper.GetString("db.timezone"),
	})
	if err != nil {
		log.Fatal("error occurred while db connection: %s", err.Error())
	}

	err = db.DB.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal("error occurred while db migration: %s", err.Error())
	}

	repos := repository.NewRepository(db.DB, log)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, log)

	srv := httpserver.New(handlers.InitRoutes(), httpserver.Port(viper.GetString("port")))

	log.Info("user-service has started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)
	<-quit

	log.Info("user-service shutting down")

	if err := srv.Shutdown(); err != nil {
		log.Error("error occurred while shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Error("error occurred while shutting down: %s", err.Error())
	}
}

func initConfigs() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
