-- name: ListUsers :many
SELECT *
FROM ibooking.users
ORDER BY user_num;

-- name: GetUser :one
SELECT *
FROM ibooking.users
WHERE user_id = $1;

-- name: GetUserByName :one
SELECT *
FROM ibooking.users
WHERE user_num = $1;


-- name: DeleteUsers :exec
DELETE
FROM ibooking.users
WHERE user_id = $1;


-- name: CreateUsers :one
INSERT INTO ibooking.users (User_Num, Password_Hash, User_Name)
VALUES ($1,
        $2,
        $3) RETURNING *;
