/* name: GetTask :one */
SELECT 
t.* 
FROM tasks as t
INNER JOIN plans as p ON t.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId);

/* name: GetTasksByDate :many */
SELECT
  t.*
FROM
  tasks AS t
WHERE
  t.date = @date
  AND t.plan_id IN (
    SELECT
      p.id
    FROM
      plans AS p
      LEFT OUTER JOIN plan_access AS pa ON p.id = pa.plan_id
    WHERE
      p.id = @planId
      AND (
        p.user = @userId OR pa.user = @userId
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
      p.id = @planId
      AND (
        p.user = @userId OR pa.user = @userId
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

/* name: CreateTaskFromTemplate :exec */
INSERT INTO tasks (
    plan_id,
    title,
    date,
    subtitle,
    description
) 
SELECT t.plan_id, t.title, @date, t.subtitle, t.description
FROM templates as t
WHERE t.id IN (SELECT t2.id FROM templates as t2
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t2.id = @templateId AND (p.user = @userId OR pa.user = @userId));

/* name: SetIsCompleteTask :exec */
UPDATE tasks 
SET is_complete = @isComplete
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: UpdateTask :exec */
UPDATE tasks 
SET title = @title, subtitle = @subtitle, date = @date, description = @description
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: DeleteTask :exec */
DELETE FROM tasks
WHERE id IN (SELECT t.id FROM tasks as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));
