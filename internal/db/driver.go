package db

import (
	"fmt"

	"github.com/Unfield/FileHopper/data"
	"github.com/Unfield/FileHopper/internal/db/drivers"
)

type DBDriver interface {
	Init(dsn string) error
	Close() error

	CreateUser(u data.User) error
	GetUser(username string) (*data.User, error)
	ListUsers() ([]data.User, error)
	UpdateUser(u data.User) error
	DeleteUser(username string) error
}

func LoadDriver(driver string) (DBDriver, error) {
	switch driver {
	case "sqlite":
		db, err := drivers.NewSqliteDriver()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Database Driver: %w", err)
		}
		return db, nil
	case "mysql":
		return nil, fmt.Errorf("not yet implemented: mysql")
	case "postgres":
		return nil, fmt.Errorf("not yet implemented: postgres")
	default:
		return nil, fmt.Errorf("invalid Database Driver '%s'", driver)
	}
}
