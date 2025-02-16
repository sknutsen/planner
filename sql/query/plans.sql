/* name: GetPlan :one */
SELECT 
* 
FROM plans
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: ListPlans :many */
SELECT 
* 
FROM plans
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.user = @userId OR pa.user = @userId);

/* name: CreatePlan :exec */
INSERT INTO plans (
    name,
    user
) VALUES (
    @name,
    @user
);

/* name: UpdatePlan :exec */
UPDATE plans 
SET name = @name
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.id = @id AND (p.user = @userId OR pa.user = @userId));

/* name: DeletePlan :exec */
DELETE FROM plans
WHERE id IN (SELECT p.id FROM plans as p
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE p.id = @id AND (p.user = @userId OR pa.user = @userId));
