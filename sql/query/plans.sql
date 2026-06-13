/* name: GetPlan :one */
SELECT 
p.id, p.name, p.user, p.updated_at, p.deleted_at 
FROM plans p
WHERE p.deleted_at IS NULL AND p.id IN (SELECT p2.id FROM plans as p2
             LEFT OUTER JOIN plan_access as pa ON p2.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE p2.id = @id AND (p2.user = @userId OR pa.user = @userId) AND p2.deleted_at IS NULL);

/* name: ListPlans :many */
SELECT DISTINCT
p.id, p.name, p.user, p.updated_at, p.deleted_at
FROM plans p
LEFT OUTER JOIN plan_access pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
WHERE (p.user = @userId OR pa.user = @userId)
AND p.deleted_at IS NULL;

/* name: ListPlansSync :many */
SELECT DISTINCT
p.id, p.name, p.user, p.updated_at, p.deleted_at
FROM plans p
LEFT OUTER JOIN plan_access pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
WHERE (p.user = @user_id OR pa.user = @user_id)
AND p.deleted_at IS NULL
AND (COALESCE(@updated_since, '') = '' OR p.updated_at >= @updated_since)
AND (COALESCE(@cursor_ts, '') = '' OR (p.updated_at > @cursor_ts OR (p.updated_at = @cursor_ts AND p.id > @cursor_id)))
ORDER BY p.updated_at ASC, p.id ASC
LIMIT @limit_count;

/* name: CreatePlan :one */
INSERT INTO plans (
    name,
    user,
    updated_at
) VALUES (
    @name,
    @user,
    strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
)
RETURNING id, name, user, updated_at, deleted_at;

/* name: UpdatePlan :exec */
UPDATE plans 
SET name = @name,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE plans.deleted_at IS NULL AND plans.id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE p.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: UpdatePlanIfMatch :execrows */
UPDATE plans 
SET name = @name,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE plans.deleted_at IS NULL AND plans.id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE p.id = @id AND (p.user = @userId OR pa.user = @userId))
AND plans.updated_at = @base_updated_at;

/* name: DeletePlan :exec */
UPDATE plans
SET deleted_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE plans.deleted_at IS NULL AND plans.id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id AND pa.deleted_at IS NULL
             WHERE p.id = @id AND (p.user = @userId OR pa.user = @userId));
