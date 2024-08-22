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
FROM
  tasks AS t
WHERE
  t.date = ?
  AND t.plan_id IN (
    SELECT
      p.id
    FROM
      plans AS p
      LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id
    WHERE
      p.id = ?
      AND (
        p.user = ?
        OR pa.user = ?
      )
  );

/* name: GetTasksByPlan :many */
SELECT
  t.*
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
      p.id = ?
      AND (
        p.user = ?
        OR pa.user = ?
      )
  )
ORDER BY t.date ASC;

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

/* name: SetIsCompleteTask :exec */
UPDATE tasks 
SET is_complete = ?
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = ? AND (p.user = ? OR pa.user = ?));

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
