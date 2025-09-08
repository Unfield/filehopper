package main

import (
	"log"

	"github.com/Unfield/FileHopper/internal/db"
	"github.com/Unfield/cascade"
)

type Config struct {
	Database struct {
		Driver string `yaml:"driver" toml:"driver" env:"DATABASE_DRIVER" flag:"database-driver"`
		DSN    string `yaml:"dsn" toml:"dsn" env:"DATABASE_DSN" flag:"database-dsn"`
	}
}

func main() {
	cfg := Config{}

	loader := cascade.NewLoader(
		cascade.WithFile("config.toml"),
		cascade.WithEnvPrefix("FILEHOPPER"),
		cascade.WithFlags(),
	)

	if err := loader.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	dbDriver, err := db.LoadDriver(cfg.Database.Driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbDriver.Init(cfg.Database.DSN); err != nil {
		log.Fatalf("failed to initialize Database Driver: %v", err)
	}

	defer dbDriver.Close()
}
