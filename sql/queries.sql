-- name: GetDebt :one
SELECT * FROM debts
WHERE user_id = $1 LIMIT 1;

-- name: AddDebt :one
INSERT INTO debts (
    amount, description, user_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: AddPlayer :one
INSERT INTO players (
    name
) VALUES (
    $1
) RETURNING *;

-- name: GetIdOfPlayer :one
SELECT id FROM players
WHERE name = $1 LIMIT 1;