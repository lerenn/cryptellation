package main

import (
	"strconv"

	"github.com/lerenn/cryptellation/v1/configs/sql"
	"github.com/lerenn/cryptellation/v1/pkg/migrator"
	"github.com/spf13/cobra"
)

var (
	mig *migrator.Migrator
)

var migrationsCmd = &cobra.Command{
	Use:     "migrations",
	Aliases: []string{"i"},
	Short:   "Manage database migrations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := callPersistentPreRunE(cmd, args); err != nil {
			return err
		}

		// Create a migrator client
		mig, err = migrator.NewMigrator(cmd.Context(), db, sql.Migrations, nil)
		return err
	},
}

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Migrate the database",
	RunE: func(cmd *cobra.Command, args []string) error {
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

func addMigrationsCommands(cmd *cobra.Command) {
	migrationsCmd.AddCommand(migrateCmd)
	migrationsCmd.AddCommand(rollbackCmd)
	cmd.AddCommand(migrationsCmd)
}
