/* name: ListPlanAccess :many */
SELECT pa.id, pa.plan_id, pa.user, pa.updated_at, pa.deleted_at FROM plan_access as pa
INNER JOIN plans as p ON pa.plan_id = p.id AND p.deleted_at IS NULL
WHERE p.id = ? AND p.user = ? AND pa.deleted_at IS NULL;

/* name: ListPlanAccessByPlanSync :many */
SELECT pa.id, pa.plan_id, pa.user, pa.updated_at, pa.deleted_at
FROM plan_access AS pa
INNER JOIN plans AS p ON pa.plan_id = p.id AND p.deleted_at IS NULL
WHERE pa.plan_id = @plan_id
  AND pa.deleted_at IS NULL
  AND (
    p.user = @user_id
    OR EXISTS (
      SELECT 1 FROM plan_access AS px
      WHERE px.plan_id = p.id AND px.user = @user_id AND px.deleted_at IS NULL
    )
  )
AND (@updated_since = '' OR pa.updated_at >= @updated_since)
AND (@cursor_ts = '' OR (pa.updated_at > @cursor_ts OR (pa.updated_at = @cursor_ts AND pa.id > @cursor_id)))
ORDER BY pa.updated_at ASC, pa.id ASC
LIMIT @limit_count;

/* name: GrantAccess :exec */
INSERT INTO plan_access (
    plan_id,
    user,
    updated_at
) VALUES (
    ?,
    ?,
    strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
);

/* name: RemoveAccess :exec */
UPDATE plan_access
SET deleted_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE user = ? AND deleted_at IS NULL;
