package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  bool
	TimeZone string
}

func (c *Config) getSSLMode() string {
	if c.SSLMode {
		return "enable"
	}
	return "disable"
}

func NewPostgresDB(config *Config) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
			config.Username, config.Password, config.DBName, config.Host, config.Port, config.getSSLMode(), config.TimeZone),
	}), &gorm.Config{})
}
