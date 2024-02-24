-- name: GetDebt :one
SELECT * FROM debt
WHERE user_id = $1 LIMIT 1;

-- name: SetDebt :one
INSERT INTO debt (
    amount, user_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: AddDebtJournalEntry :one
INSERT INTO debt_journal (
    amount, description, user_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: AddPlayer :one
INSERT INTO player (
    name
) VALUES (
    $1
) RETURNING *;

-- name: AllPlayerDebts :many
SELECT * FROM player p
JOIN debt d ON p.id = d.user_id;

-- name: GetIdOfPlayer :one
SELECT id FROM player
WHERE name = $1 LIMIT 1;

