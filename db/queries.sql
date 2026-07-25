-- name: CreateURL :one
INSERT INTO urls (code, original_url) 
VALUES ($1, $2)
RETURNING *;

-- name: GetURL :one
SELECT original_url FROM urls 
WHERE code = $1 LIMIT 1;
