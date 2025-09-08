package main

import (
	"fmt"
	"log"

	"github.com/Unfield/FileHopper/internal/auth"
	"github.com/Unfield/FileHopper/internal/db"
	"github.com/Unfield/FileHopper/internal/sftp"
	"github.com/Unfield/cascade"
)

type Config struct {
	SFTP struct {
		Port     uint32 `yaml:"port" toml:"port" env:"SFTP_PORT" flag:"sftp-port"`
		Hostname string `yaml:"hostname" toml:"hostname" env:"SFTP_HOSTNAME" flag:"sftp-hostname"`
	}
	RootDir  string `yaml:"root_dir" toml:"root_dir" env:"ROOT_DIR" flag:"root-dir"`
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

	authenticator, err := auth.NewAuthenticator(dbDriver)
	if err != nil {
		panic(err)
	}

	sftpServer := sftp.NewSFTPServer(fmt.Sprintf("%s:%d", cfg.SFTP.Hostname, cfg.SFTP.Port), authenticator, &dbDriver, cfg.RootDir)
	err = sftpServer.Start()
	if err != nil {
		panic(err)
	}
}
