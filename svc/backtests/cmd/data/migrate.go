package main

import (
	"context"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lerenn/cryptellation/pkg/adapters/db/sql"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/backtests/deployments"
	"github.com/lerenn/cryptellation/svc/backtests/internal/adapters/db/sql/migrations"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	hostFlag     string
	userFlag     string
	passwordFlag string
	databaseFlag string
	portFlag     int

	migrator *gormigrate.Gormigrate
)

var migrationsCmd = &cobra.Command{
	Use:     "migrations",
	Aliases: []string{"m"},
	Short:   "Manage migrations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		// Set default config and load from environment
		c := config.LoadSQL(&config.SQL{
			Host:     hostFlag,
			User:     userFlag,
			Password: passwordFlag,
			Database: databaseFlag,
			Port:     portFlag,
		})
		if err := c.Validate(); err != nil {
			return err
		}

		// Connect to database
		db, err := gorm.Open(postgres.Open(c.URL()), config.DefaultGormConfig)
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
		telemetry.L(cmd.Context()).Info("Launching migrations...")
		return sql.ExecuteUntilDBReady(context.TODO(), migrator.Migrate)
	},
}

func addCommandsToMigrationsCmd() {
	migrateCmd.Flags().StringVarP(&hostFlag, "host", "H", "localhost", "Host of the database")
	migrateCmd.Flags().StringVarP(&userFlag, "user", "u", "postgres", "User of the database")
	migrateCmd.Flags().StringVarP(&passwordFlag, "password", "P", "example", "Password of the database")
	migrateCmd.Flags().StringVarP(&databaseFlag, "database", "d", "backtests", "Database to migrate")
	migrateCmd.Flags().IntVarP(&portFlag, "port", "p", deployments.DockerComposeSQLDBPort, "Port of the database")

	migrationsCmd.AddCommand(migrateCmd)
}
