package migrator

import (
	"embed"
)

//go:embed migration/*.sql
var Migrations embed.FS
