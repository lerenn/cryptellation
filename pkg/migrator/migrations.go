package migrator

import "time"

// Migration is a struct that contains the migration.
type Migration struct {
	ID          int
	Description string
	SQL         string
}

// AppliedMigration is a struct that contains the status of a migration.
type AppliedMigration struct {
	ID          int       `db:"id"`
	Description string    `db:"description"`
	AppliedAt   time.Time `db:"applied_at"`
}
