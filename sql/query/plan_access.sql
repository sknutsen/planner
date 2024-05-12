/* name: ListPlanAccess :many */
SELECT pa.* FROM plan_access as pa
INNER JOIN plans as p ON pa.plan_id = p.id
WHERE p.id = ? AND p.user = ?;

/* name: GrantAccess :exec */
INSERT INTO plan_access (
    plan_id,
    user
) VALUES (
    ?,
    ?
);

/* name: RemoveAccess :exec */
DELETE FROM plan_access
WHERE user = ?;
