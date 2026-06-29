-- name: GetFeedFromURL :one
SELECT * FROM feeds
where url = $1;