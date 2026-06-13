/* name: GetResource :one */
SELECT 
    r.id, r.title, r.resource_type, r.content, r.plan_id, r.updated_at, r.deleted_at
FROM resources as r
INNER JOIN plans as p ON r.plan_id = p.id AND p.deleted_at IS NULL
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
WHERE r.id = @id AND r.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId);

/* name: GetResourcesByPlan :many */
SELECT
  r.id, r.title, r.resource_type, r.content, r.plan_id, r.updated_at, r.deleted_at
FROM
  resources AS r
WHERE
  r.deleted_at IS NULL
  AND r.plan_id IN (
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
ORDER BY r.title ASC;

/* name: ListResourcesByPlanSync :many */
SELECT
  r.id, r.title, r.resource_type, r.content, r.plan_id, r.updated_at, r.deleted_at
FROM
  resources AS r
WHERE
  r.deleted_at IS NULL
  AND r.plan_id = @plan_id
  AND r.plan_id IN (
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
AND (COALESCE(@updated_since, '') = '' OR r.updated_at >= @updated_since)
AND (COALESCE(@cursor_ts, '') = '' OR (r.updated_at > @cursor_ts OR (r.updated_at = @cursor_ts AND r.id > @cursor_id)))
ORDER BY r.updated_at ASC, r.id ASC
LIMIT @limit_count;

/* name: CreateResource :one */
INSERT INTO resources (
    plan_id,
    title,
    resource_type,
    content,
    updated_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
)
RETURNING id, title, resource_type, content, plan_id, updated_at, deleted_at;

/* name: UpdateResource :exec */
UPDATE resources 
SET title = @title, resource_type = @resourceType, content = @content,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id IN (SELECT r.id FROM resources as r
             INNER JOIN plans as p ON r.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE r.id = @id AND r.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId));

/* name: UpdateResourceIfMatch :execrows */
UPDATE resources 
SET title = @title, resource_type = @resourceType, content = @content,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id IN (SELECT r.id FROM resources as r
             INNER JOIN plans as p ON r.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE r.id = @id AND r.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId))
AND resources.updated_at = @base_updated_at;

/* name: DeleteResource :exec */
UPDATE resources
SET deleted_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id IN (SELECT r.id FROM resources as r
             INNER JOIN plans as p ON r.plan_id = p.id AND p.deleted_at IS NULL
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE r.id = @id AND r.deleted_at IS NULL AND (p.user = @userId OR pa.user = @userId));
