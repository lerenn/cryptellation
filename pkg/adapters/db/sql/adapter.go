package sql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	Client *gorm.DB
}

func New(c config.SQL) (*Adapter, error) {
	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Generate SQL client
	client, err := gorm.Open(postgres.Open(c.URL()), config.DefaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening sqldb connection: %w", err)
	}

	// Return client
	return &Adapter{
		Client: client,
	}, nil
}

func Reset(ctx context.Context, client *gorm.DB, entities []interface{}) error {
	for _, entity := range entities {
		err := client.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(entity).Error
		if err != nil {
			return fmt.Errorf("emptying %T table: %w", entity, err)
		}
	}

	return nil
}

func ExecuteUntilDBReady(ctx context.Context, fn func() error) error {
	for i := 0; i < 10; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		// SQLSTATE 3D000 = Database doesn't exist
		// SQLSTATE 42P01 = Relation doesn't exist
		// SQLSTATE 55000 = Table is being created
		// SQLSTATE 42602 = No database or schema specified
		if strings.Contains(err.Error(), "(SQLSTATE 3D000)") ||
			strings.Contains(err.Error(), "(SQLSTATE 42P01)") ||
			strings.Contains(err.Error(), "(SQLSTATE 55000)") ||
			strings.Contains(err.Error(), "(SQLSTATE 42602)") {

			telemetry.L(ctx).Debug("Database does not exist (yet). Waiting for creation and retry...")
			time.Sleep(time.Second)
			continue
		}

		return err
	}

	return fmt.Errorf("waited to long for DB to be ready")
}
