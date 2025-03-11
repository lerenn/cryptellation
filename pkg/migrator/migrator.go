package migrator

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lerenn/cryptellation/v1/pkg/telemetry"
)

const (
	migrationTableName = "_migrations"

	upMigrationKeyword   = "up"
	downMigrationKeyword = "down"
)

// Options represents the options for the migrator.
type Options struct {
	Log *log.Logger
}

// Migrator is a struct that contains all the methods to interact with the
// migrations table in the database.
type Migrator struct {
	db             *sqlx.DB
	upMigrations   []Migration
	downMigrations []Migration
}

// NewMigrator creates a new migrator.
func NewMigrator(ctx context.Context, db *sqlx.DB, upMigrations, downMigrations embed.FS, opts *Options) (*Migrator, error) {
	up, err := loadMigrations(upMigrations)
	if err != nil {
		return nil, err
	}
	telemetry.L(ctx).Infof("Loaded %d up migrations", len(up))

	down, err := loadMigrations(downMigrations)
	if err != nil {
		return nil, err
	}
	telemetry.L(ctx).Infof("Loaded %d down migrations", len(down))

	if err := setupMigrationTable(ctx, db); err != nil {
		return nil, err
	}

	return &Migrator{
		db:             db,
		upMigrations:   up,
		downMigrations: down,
	}, nil
}

func loadMigrations(migrationsDir embed.FS) ([]Migration, error) {
	entries, err := migrationsDir.ReadDir(".")
	if err != nil {
		return nil, err
	}

	migrations := make([]Migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Load migration
		parts := strings.Split(entry.Name(), ".")
		if len(parts) != 3 || parts[2] != "sql" {
			continue
		}

		// Parse migration ID
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		// Read content
		content, err := migrationsDir.ReadFile(entry.Name())
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			ID:          id,
			Description: parts[1],
			SQL:         string(content),
		})
	}

	return migrations, nil
}

func setupMigrationTable(ctx context.Context, db *sqlx.DB) error {
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS `+migrationTableName+` (
			id BIGINT PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	return nil
}

// GetLastMigrationID returns the last migration ID.
func (m *Migrator) GetLastMigrationID(ctx context.Context) (int, error) {
	var id int
	row := m.db.QueryRowContext(ctx, "SELECT id FROM "+migrationTableName+" ORDER BY id DESC LIMIT 1")
	err := row.Scan(&id)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to get last migration ID: %w", err)
	}

	return id, nil
}

// MigrateTo will apply the migrations up to the specified ID.
func (m *Migrator) MigrateTo(ctx context.Context, id int) error {
	lastID, err := m.GetLastMigrationID(ctx)
	if err != nil {
		return err
	}

	if id <= lastID {
		return errors.New("migration ID must be greater than the last migration ID")
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, migration := range m.upMigrations {
		if migration.ID <= lastID {
			continue
		}

		if migration.ID > id {
			break
		}

		if err := m.applyMigration(ctx, tx, migration, upMigrationKeyword); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *Migrator) applyMigration(ctx context.Context, tx *sql.Tx, migration Migration, direction string) error {
	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return fmt.Errorf("failed to apply migration %d: %w", migration.ID, err)
	}

	var action string
	switch direction {
	case downMigrationKeyword:
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+migrationTableName+" WHERE id = $1", migration.ID); err != nil {
			return fmt.Errorf("failed to delete migration record: %w", err)
		}
		action = "Removed"
	case upMigrationKeyword:
		_, err := tx.ExecContext(ctx,
			"INSERT INTO "+migrationTableName+" (id, description) VALUES ($1, $2)",
			migration.ID, migration.Description)
		if err != nil {
			return fmt.Errorf("failed to insert migration record: %w", err)
		}
		action = "Applied"
	default:
		return errors.New("invalid migration direction")
	}

	telemetry.L(ctx).Infof("%s %d: %s", action, migration.ID, migration.Description)
	return nil
}

// MigrateToLatest will apply all the migrations that have not been applied yet.
func (m *Migrator) MigrateToLatest(ctx context.Context) error {
	lastID, err := m.GetLastMigrationID(ctx)
	if err != nil {
		return err
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, migration := range m.upMigrations {
		if migration.ID <= lastID {
			continue
		}

		if err := m.applyMigration(ctx, tx, migration, upMigrationKeyword); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RollbackUntil will rollback the migrations until the specified ID.
func (m *Migrator) RollbackUntil(ctx context.Context, id int) error {
	lastID, err := m.GetLastMigrationID(ctx)
	if err != nil {
		return err
	}

	if id > lastID {
		return fmt.Errorf("migration ID (%d) must be less than the last migration ID (%d)",
			id, lastID)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for i := len(m.downMigrations) - 1; i >= 0; i-- {
		migration := m.downMigrations[i]
		if migration.ID > lastID {
			continue
		}

		if migration.ID < id {
			break
		}

		if err := m.applyMigration(ctx, tx, migration, downMigrationKeyword); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Rollback will rollback the last migration.
func (m *Migrator) Rollback(ctx context.Context) error {
	lastID, err := m.GetLastMigrationID(ctx)
	if err != nil {
		return err
	}

	if lastID == 0 {
		return errors.New("no migrations to rollback")
	}

	return m.RollbackUntil(ctx, lastID-1)
}
