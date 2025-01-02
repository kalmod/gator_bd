-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id) VALUES (
  $1, $2, $3, $4, $5, $6
  ) RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds set last_fetched_at = $1, updated_at = $2 WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;