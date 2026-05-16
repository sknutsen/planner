-- +goose Up
-- +goose StatementBegin
ALTER TABLE idempotency_keys ADD COLUMN request_hash TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SQLite cannot DROP COLUMN in all builds; leave column if downgrade
-- +goose StatementEnd
