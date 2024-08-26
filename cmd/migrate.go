package cmd

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"os"
	"surge/internal/conf"
	. "surge/internal/utilities"
)

var migrateCommand = cobra.Command{
	Use:     "migrate",
	Short:   "Migrate database",
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleMigrateCommand(cmd, args); err != nil {
			panic(err)
		}
	},
}

func buildMigrateCommand() *cobra.Command {
	return &migrateCommand
}

func handleMigrateCommand(cmd *cobra.Command, args []string) error {
	config, err := conf.LoadFromEnvironments()
	if err != nil {
		logrus.Errorf("Failed to load configuration from environments")
		return err
	}

	databaseUrl := config.Database.Url
	sourceUrl := StringDefault(os.Getenv("SURGE_MIGRATIONS_PATH_OVERRIDE"), "file://./migrations")

	log.Printf("%s -> %s", sourceUrl, databaseUrl)

	m, err := migrate.New(sourceUrl, databaseUrl)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("[Migrate] Nothing changed")
		} else if errors.Is(err, os.ErrNotExist) {
			log.Println("[Migrate] Not exists.")
		} else {
			return err
		}
	}

	log.Println("[Migrate] Migrate complete.")

	return nil
}
