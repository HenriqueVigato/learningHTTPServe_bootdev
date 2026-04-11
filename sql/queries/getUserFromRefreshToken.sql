-- name: GetUserFromRefresToken :one
SElECT users.* FROM refresh_tokens
  INNER JOIN users on refresh_tokens.user_id = users.id
  WHERE  refresh_tokens.token = $1;
