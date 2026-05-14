package handler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

// OpenTursoDB opens a libsql client for reuse across requests. The caller owns
// the returned *sql.DB and must Close it on process shutdown.
func OpenTursoDB(cfg TursoConfig) (*sql.DB, error) {
	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", cfg.PrimaryUrl, cfg.AuthToken))
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(9 * time.Second)
	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
