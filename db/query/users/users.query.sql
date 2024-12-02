-- name: GetUserByID :one
SELECT 
    id, 
    username, 
    email, 
    phone, 
    status, 
    preference_low_channel, 
    preference_medium_channel, 
    preference_high_channel 
FROM 
    users 
WHERE 
    id = $1;

-- name: CreateUser :exec
INSERT INTO users (
    id, 
    username, 
    email, 
    phone, 
    status, 
    preference_low_channel, 
    preference_medium_channel, 
    preference_high_channel
) 
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: DeleteUserByID :exec
DELETE FROM users 
WHERE id = $1;


-- name: GetAllUsers :many
SELECT 
    id, 
    username, 
    email, 
    phone, 
    status, 
    preference_low_channel, 
    preference_medium_channel, 
    preference_high_channel
FROM 
    users;
