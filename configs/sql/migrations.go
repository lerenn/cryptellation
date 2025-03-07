package sql

import "embed"

// Migrations contains all the migrations to be applied to the database.
//
//go:embed *.sql
var Migrations embed.FS
