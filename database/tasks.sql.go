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

const createTaskFromTemplate = `-- name: CreateTaskFromTemplate :exec
INSERT INTO tasks (
    plan_id,
    title,
    date,
    subtitle,
    description
) 
SELECT t.plan_id, t.title, ?1, t.subtitle, t.description
FROM templates as t
WHERE t.id IN (SELECT t2.id FROM templates as t2
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t2.id = ?2 AND (p.user = ?3 OR pa.user = ?3))
`

type CreateTaskFromTemplateParams struct {
	Date       string
	TemplateId int64
	UserId     string
}

func (q *Queries) CreateTaskFromTemplate(ctx context.Context, arg CreateTaskFromTemplateParams) error {
	_, err := q.db.ExecContext(ctx, createTaskFromTemplate, arg.Date, arg.TemplateId, arg.UserId)
	return err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ?1 AND (p.user = ?2 OR pa.user = ?2))
`

type DeleteTaskParams struct {
	ID     int64
	UserId string
}

func (q *Queries) DeleteTask(ctx context.Context, arg DeleteTaskParams) error {
	_, err := q.db.ExecContext(ctx, deleteTask, arg.ID, arg.UserId)
	return err
}

const getTask = `-- name: GetTask :one
SELECT 
t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete 
FROM tasks as t
INNER JOIN plans as p ON t.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE t.id = ?1 AND (p.user = ?2 OR pa.user = ?2)
`

type GetTaskParams struct {
	ID     int64
	UserId string
}

func (q *Queries) GetTask(ctx context.Context, arg GetTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTask, arg.ID, arg.UserId)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.PlanID,
		&i.Date,
		&i.Title,
		&i.Subtitle,
		&i.Description,
		&i.IsComplete,
	)
	return i, err
}

const getTasksByDate = `-- name: GetTasksByDate :many
SELECT
  t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete
FROM
  tasks AS t
WHERE
  t.date = ?1
  AND t.plan_id IN (
    SELECT
      p.id
    FROM
      plans AS p
      LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id
    WHERE
      p.id = ?2
      AND (
        p.user = ?3 OR pa.user = ?3
      )
  )
`

type GetTasksByDateParams struct {
	Date   string
	PlanId int64
	UserId string
}

func (q *Queries) GetTasksByDate(ctx context.Context, arg GetTasksByDateParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksByDate, arg.Date, arg.PlanId, arg.UserId)
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
			&i.IsComplete,
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

const getTasksByPlan = `-- name: GetTasksByPlan :many
SELECT
  t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete
FROM
  tasks AS t
WHERE
  t.plan_id IN (
    SELECT
      p.id
    FROM
      plans AS p
      LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id
    WHERE
      p.id = ?1
      AND (
        p.user = ?2 OR pa.user = ?2
      )
  )
ORDER BY t.date ASC
`

type GetTasksByPlanParams struct {
	PlanId int64
	UserId string
}

func (q *Queries) GetTasksByPlan(ctx context.Context, arg GetTasksByPlanParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksByPlan, arg.PlanId, arg.UserId)
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
			&i.IsComplete,
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

const setIsCompleteTask = `-- name: SetIsCompleteTask :exec
UPDATE tasks 
SET is_complete = ?1
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ?2 AND (p.user = ?3 OR pa.user = ?3))
`

type SetIsCompleteTaskParams struct {
	IsComplete int64
	ID         int64
	UserId     string
}

func (q *Queries) SetIsCompleteTask(ctx context.Context, arg SetIsCompleteTaskParams) error {
	_, err := q.db.ExecContext(ctx, setIsCompleteTask, arg.IsComplete, arg.ID, arg.UserId)
	return err
}

const updateTask = `-- name: UpdateTask :exec
UPDATE tasks 
SET title = ?1, subtitle = ?2, date = ?3, description = ?4
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ?5 AND (p.user = ?6 OR pa.user = ?6))
`

type UpdateTaskParams struct {
	Title       string
	Subtitle    interface{}
	Date        string
	Description interface{}
	ID          int64
	UserId      string
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) error {
	_, err := q.db.ExecContext(ctx, updateTask,
		arg.Title,
		arg.Subtitle,
		arg.Date,
		arg.Description,
		arg.ID,
		arg.UserId,
	)
	return err
}
