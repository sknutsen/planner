-- +goose Up
CREATE TABLE plan_access (
    id integer PRIMARY KEY AUTOINCREMENT,
    plan_id integer NOT NULL,
    user text NOT NULL
);

-- +goose Down
DROP TABLE plan_access;
