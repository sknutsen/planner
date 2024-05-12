package handler

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func (h *Handler) openDB() *sql.DB {
	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", h.TursoConfig.PrimaryUrl, h.TursoConfig.AuthToken))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed creating authenticator. Error: %s\n", err)
		os.Exit(1)
	}

	return db
}
