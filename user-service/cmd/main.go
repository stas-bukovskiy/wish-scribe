package main

import (
	"context"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	userService "github.com/stas-bukovskiy/wish-scribe/user-service"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/entity"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/handler"
	repository2 "github.com/stas-bukovskiy/wish-scribe/user-service/internal/repository"
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
	if err := initConfigs(); err != nil {
		log.Fatalf("error occurred while config initalizing: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error occurred while env variables loading: %s", err.Error())
	}

	db, err := repository2.NewPostgresDB(&repository2.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetBool("db.sslmode"),
		TimeZone: viper.GetString("db.timezone"),
	})
	if err != nil {
		log.Fatalf("error occurred while db connection: %s", err.Error())
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("error occurred while db migration: %s", err.Error())
	}

	repos := repository2.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(userService.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occurred while running http httpserver: %s", err.Error())
		}
	}()

	log.Printf("user-service has started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)
	<-quit

	log.Printf("user-service shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occurred while shutting down: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("error occurred while getting sql db: %s", err.Error())
	}
	if err := sqlDB.Close(); err != nil {
		log.Errorf("error occurred while shutting down: %s", err.Error())
	}
}

func initConfigs() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
