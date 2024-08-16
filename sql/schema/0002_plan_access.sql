-- +goose Up
-- +goose StatementBegin
CREATE TABLE plan_access (
    id integer PRIMARY KEY AUTOINCREMENT,
    plan_id integer NOT NULL,
    user text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE plan_access;
-- +goose StatementEnd
