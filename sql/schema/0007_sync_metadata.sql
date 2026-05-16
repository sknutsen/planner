-- +goose Up
-- SQLite (incl. Turso) rejects non-constant DEFAULT on ALTER TABLE ADD COLUMN; use a literal, then backfill.
-- +goose StatementBegin
ALTER TABLE plans ADD COLUMN updated_at TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z';
UPDATE plans SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ','now');
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE plans ADD COLUMN deleted_at TEXT;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE tasks ADD COLUMN updated_at TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z';
UPDATE tasks SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ','now');
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE tasks ADD COLUMN deleted_at TEXT;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE resources ADD COLUMN updated_at TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z';
UPDATE resources SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ','now');
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE resources ADD COLUMN deleted_at TEXT;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE templates ADD COLUMN updated_at TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z';
UPDATE templates SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ','now');
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE templates ADD COLUMN deleted_at TEXT;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE plan_access ADD COLUMN updated_at TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z';
UPDATE plan_access SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ','now');
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE idempotency_keys (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    key_hash TEXT NOT NULL,
    response_body TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now')),
    UNIQUE(user_id, key_hash)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS idempotency_keys;
-- +goose StatementEnd
