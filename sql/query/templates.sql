/* name: GetTemplate :one */
SELECT 
t.* 
FROM templates as t
INNER JOIN plans as p ON t.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId);

/* name: GetTemplatesByPlan :many */
SELECT
  t.*
FROM
  templates AS t
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
ORDER BY t.title ASC;

/* name: CreateTemplate :exec */
INSERT INTO templates (
    plan_id,
    title,
    subtitle,
    description
) VALUES (
    ?,
    ?,
    ?,
    ?
);

/* name: UpdateTemplate :exec */
UPDATE templates 
SET title = @title, subtitle = @subtitle, description = @description
WHERE id IN (SELECT t.id FROM templates as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: DeleteTemplate :exec */
DELETE FROM templates
WHERE id IN (SELECT t.id FROM templates as t
             INNER JOIN plans as p ON t.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE t.id = @id AND (p.user = @userId OR pa.user = @userId));
