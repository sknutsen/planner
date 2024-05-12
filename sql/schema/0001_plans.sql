-- +goose Up
CREATE TABLE plans (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL,
    user text NOT NULL
);

-- +goose Down
DROP TABLE plans;
