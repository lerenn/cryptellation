package up

import "embed"

// Migrations contains all the migrations to be applied to the database when
// migrating.
//
//go:embed *.sql
var Migrations embed.FS
