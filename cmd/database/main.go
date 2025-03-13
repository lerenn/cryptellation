package main

import (
	"context"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/console"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry/otel"
	"github.com/lerenn/cryptellation/v1/pkg/version"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var db *sqlx.DB

var (
	driverNameFlag string
	dsnFlag        string
)

var rootCmd = &cobra.Command{
	Use:     os.Args[0],
	Version: version.FullVersion(),
	Short:   os.Args[0] + " - a CLI to manage cryptellation database",
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

func main() {
	var errCode int

	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation"))

	// Set flags
	c := config.LoadSQL(nil)
	rootCmd.PersistentFlags().StringVarP(&driverNameFlag, "driver", "d", "postgres", "Set the database driver name")
	rootCmd.PersistentFlags().StringVarP(&dsnFlag, "dsn", "s", c.DSN, "Set the database data source name")

	// Add commands
	addMigrationsCommands(rootCmd)

	// Execute command
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		telemetry.L(context.Background()).Errorf("an error occurred: %s", err.Error())
	}

	// Close telemetry
	telemetry.Close(context.Background())

	// Exit with error code
	os.Exit(errCode)
}
