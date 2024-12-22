-- +goose Up
-- +goose StatementBegin
CREATE TABLE resources (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    title text NOT NULL,
    resource_type integer NOT NULL,
    content text NULL,
    plan_id integer NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resources;
-- +goose StatementEnd
