package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := applySchemaMigrations(t, db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	return db
}

func applySchemaMigrations(t *testing.T, db *sql.DB) error {
	t.Helper()
	dir := filepath.Join("..", "sql", "schema")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		for _, stmt := range gooseUpStatements(string(data)) {
			if _, err := db.Exec(stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

func gooseUpStatements(content string) []string {
	if idx := strings.Index(content, "-- +goose Down"); idx >= 0 {
		content = content[:idx]
	}
	upIdx := strings.Index(content, "-- +goose Up")
	if upIdx < 0 {
		return nil
	}
	content = content[upIdx:]

	var stmts []string
	for {
		begin := strings.Index(content, "-- +goose StatementBegin")
		if begin < 0 {
			break
		}
		content = content[begin+len("-- +goose StatementBegin"):]
		end := strings.Index(content, "-- +goose StatementEnd")
		if end < 0 {
			break
		}
		stmt := strings.TrimSpace(content[:end])
		if stmt != "" {
			stmts = append(stmts, stmt)
		}
		content = content[end+len("-- +goose StatementEnd"):]
	}
	return stmts
}
