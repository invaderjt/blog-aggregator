-- name: GetNameFromUUID :one
SELECT name FROM users
where ID = $1;