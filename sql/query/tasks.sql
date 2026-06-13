/* name: GetTask :one */
SELECT 
t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete, t.updated_at, t.deleted_at
FROM tasks as t
INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
WHERE t.id = @id AND t.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId);

/* name: GetTasksByDate :many */
SELECT
  t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete, t.updated_at, t.deleted_at
FROM
  tasks AS t
WHERE
  t.deleted_at IS NULL
  AND t.date = @date
  AND t.plan_id IN (
    SELECT
      p.id
    FROM
      plans AS p
      LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
    WHERE
      p.id = @planId
      AND p.deleted_at IS NULL
      AND (
        p.user = @userId OR pa.user = @userId
      )
  );

/* name: GetTasksByPlan :many */
SELECT
  t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete, t.updated_at, t.deleted_at
FROM
  tasks AS t
WHERE
  t.deleted_at IS NULL
  AND t.plan_id IN (
    SELECT
      p.id
    FROM
      plans AS p
      LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
    WHERE
      p.id = @planId
      AND p.deleted_at IS NULL
      AND (
        p.user = @userId OR pa.user = @userId
      )
  )
ORDER BY t.date ASC;

/* name: ListTasksByPlanSync :many */
SELECT
  t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete, t.updated_at, t.deleted_at
FROM
  tasks AS t
WHERE
  t.deleted_at IS NULL
  AND t.plan_id = @plan_id
  AND t.plan_id IN (
    SELECT
      p.id
    FROM
      plans AS p
      LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
    WHERE
      p.id = @plan_id
      AND p.deleted_at IS NULL
      AND (
        p.user = @user_id OR pa.user = @user_id
      )
  )
AND (COALESCE(@updated_since, '') = '' OR t.updated_at >= @updated_since)
AND (COALESCE(@cursor_ts, '') = '' OR (t.updated_at > @cursor_ts OR (t.updated_at = @cursor_ts AND t.id > @cursor_id)))
ORDER BY t.updated_at ASC, t.id ASC
LIMIT @limit_count;

/* name: ListTasksByPlanAndDates :many */
SELECT
  t.id, t.plan_id, t.date, t.title, t.subtitle, t.description, t.is_complete, t.updated_at, t.deleted_at
FROM
  tasks AS t
WHERE
  t.deleted_at IS NULL
  AND t.plan_id = @plan_id
  AND EXISTS (
    SELECT 1
    FROM plans AS p
    LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
    WHERE p.id = @plan_id
      AND p.deleted_at IS NULL
      AND (p.user = @user_id OR pa.user = @user_id)
  )
  AND t.date IN (sqlc.slice('dates'))
ORDER BY t.date ASC;

/* name: CreateTask :one */
INSERT INTO tasks (
    plan_id,
    title,
    date,
    subtitle,
    description,
    updated_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
)
RETURNING id, plan_id, date, title, subtitle, description, is_complete, updated_at, deleted_at;

/* name: CreateTaskFromTemplate :one */
INSERT INTO tasks (
    plan_id,
    title,
    date,
    subtitle,
    description,
    updated_at
) 
SELECT t.plan_id, t.title, @date, t.subtitle, t.description, strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
FROM templates as t
WHERE t.deleted_at IS NULL AND t.id IN (SELECT t2.id FROM templates as t2
             INNER JOIN plans as p ON t2.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t2.id = @templateId AND t2.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId))
RETURNING id, plan_id, date, title, subtitle, description, is_complete, updated_at, deleted_at;

/* name: SetIsCompleteTask :exec */
UPDATE tasks 
SET is_complete = @isComplete,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE deleted_at IS NULL AND id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: UpdateTask :exec */
UPDATE tasks 
SET title = @title, subtitle = @subtitle, date = @date, description = @description,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE deleted_at IS NULL AND id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: UpdateTaskIfMatch :execrows */
UPDATE tasks 
SET title = @title, subtitle = @subtitle, date = @date, description = @description,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE deleted_at IS NULL AND id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId))
AND tasks.updated_at = @base_updated_at;

/* name: DeleteTask :exec */
UPDATE tasks
SET deleted_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE deleted_at IS NULL AND id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));
