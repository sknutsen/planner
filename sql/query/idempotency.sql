/* name: GetIdempotencyRecord :one */
SELECT request_hash, response_body, status_code FROM idempotency_keys
WHERE user_id = ? AND key_hash = ?;

/* name: InsertIdempotencyResponse :exec */
INSERT INTO idempotency_keys (user_id, key_hash, request_hash, response_body, status_code, created_at)
VALUES (?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ', 'now'));
