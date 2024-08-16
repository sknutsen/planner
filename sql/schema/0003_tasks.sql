-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id integer PRIMARY KEY AUTOINCREMENT,
    plan_id integer NOT NULL,
    date text NOT NULL,
    title text NOT NULL,
    subtitle text NULL,
    description text NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
