package database

import (
	"context"
	"errors"
)

type Database interface {
	// Close closes the connection to database.
	Close() error
	// Ping - checks if database is available.
	Ping(ctx context.Context) error
}

var (
	ErrDatabaseEmptyOutput = errors.New("empty query output")
)
