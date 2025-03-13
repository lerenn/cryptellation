package down

import "embed"

// Migrations contains all the migrations to be applied to the database when
// rollbacking.
//
//go:embed *.sql
var Migrations embed.FS
