/* name: GetTask :one */
SELECT 
t.* 
FROM tasks as t
INNER JOIN plans as p ON t.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE t.id = ? AND (p.user = ? OR pa.user = ?);

/* name: GetTasksByDate :many */
SELECT 
t.* 
FROM tasks as t
INNER JOIN plans as p ON t.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE t.date = ? AND (p.user = ? OR pa.user = ?);

/* name: CreateTask :exec */
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
);

/* name: UpdateTask :exec */
UPDATE tasks 
SET title = ?, subtitle = ?, date = ?, description = ?
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ? AND (p.user = ? OR pa.user = ?));

/* name: DeleteTask :exec */
DELETE FROM tasks
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ? AND (p.user = ? OR pa.user = ?));
