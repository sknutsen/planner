-- +goose Up
-- +goose StatementBegin
CREATE TABLE templates (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    plan_id integer NOT NULL,
    title text NOT NULL,
    subtitle text NULL,
    description text NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resources;
-- +goose StatementEnd
