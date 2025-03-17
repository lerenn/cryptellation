package main

import (
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lerenn/cryptellation/v1/configs/sql/down"
	"github.com/lerenn/cryptellation/v1/configs/sql/up"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/migrator"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	driverNameFlag string
	dsnFlag        string
)

var (
	db *sqlx.DB
)

var databaseCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"i"},
	Short:   "Manage database",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) (err error) {
		// Create a sqlx client
		for {
			db, err = sqlx.ConnectContext(cmd.Context(), driverNameFlag, dsnFlag)
			if err == nil {
				return nil
			}

			telemetry.L(cmd.Context()).Errorf("database connection failed: %s", err.Error())
			telemetry.L(cmd.Context()).Infof("retrying in 1 second...")
			time.Sleep(time.Second)
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Migrate the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a migrator client
		mig, err := migrator.NewMigrator(cmd.Context(), db, up.Migrations, down.Migrations, nil)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return mig.MigrateToLatest(cmd.Context())
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return mig.MigrateTo(cmd.Context(), id)
	},
}

var rollbackCmd = &cobra.Command{
	Use:     "rollback",
	Aliases: []string{"r"},
	Short:   "Rollback the databas before a migration ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a migrator client
		mig, err := migrator.NewMigrator(cmd.Context(), db, up.Migrations, down.Migrations, nil)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return mig.Rollback(cmd.Context())
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return mig.RollbackUntil(cmd.Context(), id)
	},
}

func addDatabaseCommands(cmd *cobra.Command) {
	databaseCmd.AddCommand(migrateCmd)
	databaseCmd.AddCommand(rollbackCmd)

	// Set flags
	c := config.LoadSQL(nil)
	databaseCmd.PersistentFlags().StringVarP(&driverNameFlag, "driver", "d", "postgres", "Set the database driver name")
	databaseCmd.PersistentFlags().StringVarP(&dsnFlag, "dsn", "s", c.DSN, "Set the database data source name")

	cmd.AddCommand(databaseCmd)
}
