package migrator

import "time"

type Migration struct {
	ID          int
	Description string
	Domain      string
	Direction   string
	SQL         string
}

type MigrationStatus struct {
	ID          int       `db:"id"`
	Description string    `db:"description"`
	Domain      string    `db:"domain"`
	AppliedAt   time.Time `db:"applied_at"`
}
