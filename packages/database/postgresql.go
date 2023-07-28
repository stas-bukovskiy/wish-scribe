package database

import (
	"context"
	"fmt"
	"time"

	// third party
	"github.com/a631807682/zerofield"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgreSQLConfig - represents PostgreSQL service config.
type PostgreSQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  bool
	TimeZone string
}

func (c *PostgreSQLConfig) getSSLMode() string {
	if c.SSLMode {
		return "enable"
	}
	return "disable"
}

// PostgreSQL - represents postgresql service.
type PostgreSQL struct {
	DB *gorm.DB
}

var _ Database = (*PostgreSQL)(nil)

// PostgreSQLModel provides base fields for database models (like gorm.PostgreSQLModel).
type PostgreSQLModel struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" swaggerignore:"true"`
}

// Option - represents PostgreSQL service option.
type Option func(*PostgreSQL)

// SetMaxIdleConns - configures max idle connections.
func SetMaxIdleConns(idleConns int) Option {
	return func(p *PostgreSQL) {
		db, _ := p.DB.DB()
		db.SetMaxIdleConns(idleConns)
	}
}

// SetMaxOpenConns - configures max open connections.
func SetMaxOpenConns(openConns int) Option {
	return func(p *PostgreSQL) {
		db, _ := p.DB.DB()
		db.SetMaxOpenConns(openConns)
	}
}

// SetConnMaxLifetime - configures max connection lifetime.
func SetConnMaxLifetime(maxLifetime time.Duration) Option {
	return func(p *PostgreSQL) {
		db, _ := p.DB.DB()
		db.SetConnMaxLifetime(maxLifetime)
	}
}

// NewPostgreSQL - creates new instance of PostgreSQL service.
func NewPostgreSQL(config PostgreSQLConfig, opts ...Option) (*PostgreSQL, error) {
	// create instance of postgresql
	sql := &PostgreSQL{}

	// apply custom options
	for _, opt := range opts {
		opt(sql)
	}

	// connect to database
	var err error
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		config.Username, config.Password, config.DBName, config.Host, config.Port, config.getSSLMode(), config.TimeZone)

	sql.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
	}

	// https://github.com/a631807682/zerofield
	// allow update zero value field
	err = sql.DB.Use(zerofield.NewPlugin())
	if err != nil {
		return nil, fmt.Errorf("failed to use zerofield plugin: %w", err)
	}

	// create UUID extension.
	err = sql.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	return sql, nil
}

func (p *PostgreSQL) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	pgSql, err := p.DB.WithContext(ctx).DB()
	if err != nil {
		return fmt.Errorf("failed to connect")
	}
	err = pgSql.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect")
	}
	return nil
}

// Close - closes postgresql service database connection.
func (p *PostgreSQL) Close() error {
	if p.DB != nil {
		db, _ := p.DB.DB()
		err := db.Close()
		if err != nil {
			return fmt.Errorf("failed to close postgresql connection: %w", err)
		}
	}
	return nil
}
