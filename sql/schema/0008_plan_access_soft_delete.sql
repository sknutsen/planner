-- +goose Up
-- +goose StatementBegin
ALTER TABLE plan_access ADD COLUMN deleted_at TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SQLite: column removal not supported without table rebuild; no-op down.
SELECT 1;
-- +goose StatementEnd
