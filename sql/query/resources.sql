/* name: GetResource :one */
SELECT 
    r.* 
FROM resources as r
INNER JOIN plans as p ON r.plan_id = p.id
LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
WHERE r.id = ? AND (p.user = ? OR pa.user = ?);

/* name: GetResourcesByPlan :many */
SELECT
  r.*
FROM
  resources AS r
WHERE
  r.plan_id IN (
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
ORDER BY r.title ASC;

/* name: CreateResource :exec */
INSERT INTO resources (
    plan_id,
    title,
    resource_type,
    content
) VALUES (
    ?,
    ?,
    ?,
    ?
);

/* name: UpdateResource :exec */
UPDATE resources 
SET title = ?, resource_type = ?, content = ?
WHERE id IN (SELECT r.id FROM resources as r
             INNER JOIN plans as p ON r.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE r.id = ? AND (p.user = ? OR pa.user = ?));

/* name: DeleteResource :exec */
DELETE FROM resources
WHERE id IN (SELECT r.id FROM resources as r
             INNER JOIN plans as p ON r.plan_id = p.id
             LEFT OUTER JOIN plan_access as pa ON p.id = pa.plan_id
             WHERE r.id = ? AND (p.user = ? OR pa.user = ?));
