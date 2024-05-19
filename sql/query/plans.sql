/* name: GetPlan :one */
SELECT 
* 
FROM plans
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.id = ? AND (p.user = ? OR pa.user = ?));

/* name: ListPlans :many */
SELECT 
* 
FROM plans
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.user = ? OR pa.user = ?);

/* name: CreatePlan :exec */
INSERT INTO plans (
    name,
    user
) VALUES (
    ?,
    ?
);

/* name: UpdatePlan :exec */
UPDATE plans 
SET name = ?
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.id = ? AND (p.user = ? OR pa.user = ?));

/* name: DeletePlan :exec */
DELETE FROM plans
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.id = ? AND (p.user = ? OR pa.user = ?));
