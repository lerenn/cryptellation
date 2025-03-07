package migrator

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	migrationTableName = "_migrations"
)

type MigratorOptions struct {
	Log *log.Logger
}

type Migrator struct {
	db         *sqlx.DB
	migrations []Migration
	logger     *log.Logger
}

func NewMigrator(ctx context.Context, db *sqlx.DB, migrations embed.FS, opts *MigratorOptions) (*Migrator, error) {
	migs, err := loadMigrations(migrations)
	if err != nil {
		return nil, err
	}

	if err := setupMigrationTable(ctx, db); err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	if opts != nil && opts.Log != nil {
		logger = opts.Log
	}

	return &Migrator{
		db:         db,
		migrations: migs,
		logger:     logger,
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
		if len(parts) != 5 || parts[4] != "sql" {
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
			Domain:      parts[2],
			Direction:   parts[3],
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
			domain TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	return nil
}

func (m *Migrator) GetLastMigrationID(ctx context.Context) (int, error) {
	var id int
	row := m.db.QueryRowContext(ctx, "SELECT id FROM "+migrationTableName+" ORDER BY id DESC LIMIT 1")
	err := row.Scan(&id)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to get last migration ID: %w", err)
	}

	return id, nil
}

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
	defer tx.Rollback()

	for _, migration := range m.migrations {
		if migration.ID <= lastID || migration.Direction == "down" {
			continue
		}

		if migration.ID > id {
			break
		}

		if err := m.applyMigration(ctx, tx, migration); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *Migrator) applyMigration(ctx context.Context, tx *sql.Tx, migration Migration) error {
	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return fmt.Errorf("failed to apply migration %d: %w", migration.ID, err)
	}

	var direction string
	if migration.Direction == "down" {
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+migrationTableName+" WHERE id = $1", migration.ID); err != nil {
			return fmt.Errorf("failed to delete migration record: %w", err)
		}
		direction = "---"
	} else if migration.Direction == "up" {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO "+migrationTableName+" (id, description, domain) VALUES ($1, $2, $3)",
			migration.ID, migration.Description, migration.Domain)
		if err != nil {
			return fmt.Errorf("failed to insert migration record: %w", err)
		}
		direction = "+++"
	} else {
		return errors.New("invalid migration direction")
	}

	m.logger.Printf("%s [%d] %-15s: %s\n", direction, migration.ID, migration.Domain, migration.Description)
	return nil
}

func (m *Migrator) MigrateToLatest(ctx context.Context) error {
	lastID, err := m.GetLastMigrationID(ctx)
	if err != nil {
		return err
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, migration := range m.migrations {
		if migration.ID <= lastID || migration.Direction == "down" {
			continue
		}

		if err := m.applyMigration(ctx, tx, migration); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

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
	defer tx.Rollback()

	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]
		if migration.ID > lastID || migration.Direction == "up" {
			continue
		}

		if migration.ID < id {
			break
		}

		if err := m.applyMigration(ctx, tx, migration); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
