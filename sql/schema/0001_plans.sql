-- +goose Up
-- +goose StatementBegin
CREATE TABLE plans (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL,
    user text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE plans;
-- +goose StatementEnd
