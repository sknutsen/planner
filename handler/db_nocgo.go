//go:build !cgo

package handler

import (
	"database/sql"
	"fmt"
)

// OpenTursoDB is unavailable when CGO is disabled. The Turso libsql driver
// requires CGO and a C toolchain (e.g. gcc). Enable CGO to build with libsql.
func OpenTursoDB(cfg TursoConfig) (*sql.DB, error) {
	_ = cfg
	return nil, fmt.Errorf("libsql driver requires CGO (set CGO_ENABLED=1 and install a C compiler)")
}
