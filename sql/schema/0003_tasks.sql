-- +goose Up
CREATE TABLE tasks (
    id integer PRIMARY KEY AUTOINCREMENT,
    plan_id integer NOT NULL,
    date text NOT NULL,
    title text NOT NULL,
    subtitle text NULL,
    description text NULL
);

-- +goose Down
DROP TABLE tasks;
