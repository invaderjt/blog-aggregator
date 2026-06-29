-- name: GetFeedFollowsForUser :many
SELECT feeds.name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;