// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: tasks.sql

package database

import (
	"context"
)

const createTask = `-- name: CreateTask :exec
INSERT INTO tasks (
    plan_id,
    title,
    date,
    subtitle,
    description
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
)
`

type CreateTaskParams struct {
	PlanID      int64
	Title       string
	Date        string
	Subtitle    interface{}
	Description interface{}
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) error {
	_, err := q.db.ExecContext(ctx, createTask,
		arg.PlanID,
		arg.Title,
		arg.Date,
		arg.Subtitle,
		arg.Description,
	)
	return err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ? AND (p.user = ? OR pa.user = ?))
`

type DeleteTaskParams struct {
	ID     int64
	User   string
	User_2 string
}

func (q *Queries) DeleteTask(ctx context.Context, arg DeleteTaskParams) error {
	_, err := q.db.ExecContext(ctx, deleteTask, arg.ID, arg.User, arg.User_2)
	return err
}

const getTask = `-- name: GetTask :one
SELECT 
t.id, t.plan_id, t.date, t.title, t.subtitle, t.description 
FROM tasks as t
INNER JOIN plans as p ON t.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE t.id = ? AND (p.user = ? OR pa.user = ?)
`

type GetTaskParams struct {
	ID     int64
	User   string
	User_2 string
}

func (q *Queries) GetTask(ctx context.Context, arg GetTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTask, arg.ID, arg.User, arg.User_2)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.PlanID,
		&i.Date,
		&i.Title,
		&i.Subtitle,
		&i.Description,
	)
	return i, err
}

const getTasksByDate = `-- name: GetTasksByDate :many
SELECT 
t.id, t.plan_id, t.date, t.title, t.subtitle, t.description 
FROM tasks as t
INNER JOIN plans as p ON t.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE t.date = ? AND t.plan_id = ? AND (p.user = ? OR pa.user = ?)
`

type GetTasksByDateParams struct {
	Date   string
	PlanID int64
	User   string
	User_2 string
}

func (q *Queries) GetTasksByDate(ctx context.Context, arg GetTasksByDateParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksByDate,
		arg.Date,
		arg.PlanID,
		arg.User,
		arg.User_2,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.PlanID,
			&i.Date,
			&i.Title,
			&i.Subtitle,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :exec
UPDATE tasks 
SET title = ?, subtitle = ?, date = ?, description = ?
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ? AND (p.user = ? OR pa.user = ?))
`

type UpdateTaskParams struct {
	Title       string
	Subtitle    interface{}
	Date        string
	Description interface{}
	ID          int64
	User        string
	User_2      string
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) error {
	_, err := q.db.ExecContext(ctx, updateTask,
		arg.Title,
		arg.Subtitle,
		arg.Date,
		arg.Description,
		arg.ID,
		arg.User,
		arg.User_2,
	)
	return err
}
