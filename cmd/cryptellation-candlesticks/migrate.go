package main

import (
	"log"

	"github.com/digital-feather/cryptellation/internal/candlesticks/infra/sql/migrations"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DefaultGormConfig = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	migrator *gormigrate.Gormigrate
)

var migrationsCmd = &cobra.Command{
	Use:     "migrations",
	Aliases: []string{"m"},
	Short:   "Manage migrations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Load config from environment
		c := config.LoadSQLConfigFromEnv()
		if err := c.Validate(); err != nil {
			return err
		}

		// Connect to database
		db, err := gorm.Open(postgres.Open(c.URL()), DefaultGormConfig)
		if err != nil {
			return err
		}

		// Register migrations
		migrator = gormigrate.New(db, gormigrate.DefaultOptions, migrations.Migrations)

		// Register initSchema in case of new database
		migrator.InitSchema(migrations.InitSchema)

		return nil
	},
}

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Execute migrations",
	Long:    "Execute all migrations that have not been applied.",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Launching migrations...")
		return migrator.Migrate()
	},
}

func addCommandsToMigrationsCmd() {
	migrationsCmd.AddCommand(migrateCmd)
}
