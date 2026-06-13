/* name: GetTemplate :one */
SELECT 
t.id, t.plan_id, t.title, t.subtitle, t.description, t.updated_at, t.deleted_at
FROM templates as t
INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
WHERE t.id = @id AND t.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId);

/* name: GetTemplatesByPlan :many */
SELECT
  t.id, t.plan_id, t.title, t.subtitle, t.description, t.updated_at, t.deleted_at
FROM
  templates AS t
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
      AND (
        p.user = @userId OR pa.user = @userId
      )
      AND p.deleted_at IS NULL
  )
ORDER BY t.title ASC;

/* name: ListTemplatesByPlanSync :many */
SELECT
  t.id, t.plan_id, t.title, t.subtitle, t.description, t.updated_at, t.deleted_at
FROM
  templates AS t
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
      AND (
        p.user = @user_id OR pa.user = @user_id
      )
      AND p.deleted_at IS NULL
  )
AND (COALESCE(@updated_since, '') = '' OR t.updated_at >= @updated_since)
AND (COALESCE(@cursor_ts, '') = '' OR (t.updated_at > @cursor_ts OR (t.updated_at = @cursor_ts AND t.id > @cursor_id)))
ORDER BY t.updated_at ASC, t.id ASC
LIMIT @limit_count;

/* name: CreateTemplate :one */
INSERT INTO templates (
    plan_id,
    title,
    subtitle,
    description,
    updated_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
)
RETURNING id, plan_id, title, subtitle, description, updated_at, deleted_at;

/* name: UpdateTemplate :exec */
UPDATE templates 
SET title = @title, subtitle = @subtitle, description = @description,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id IN (SELECT t.id FROM templates as t
             INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t.id = @id AND t.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId));

/* name: UpdateTemplateIfMatch :execrows */
UPDATE templates 
SET title = @title, subtitle = @subtitle, description = @description,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id IN (SELECT t.id FROM templates as t
             INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t.id = @id AND t.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId))
AND templates.updated_at = @base_updated_at;

/* name: DeleteTemplate :exec */
UPDATE templates
SET deleted_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id IN (SELECT t.id FROM templates as t
             INNER JOIN plans as p ON t.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE t.id = @id AND t.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId));
