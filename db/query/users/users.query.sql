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
