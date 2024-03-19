-- name: GetDebt :one
SELECT * FROM debt
WHERE user_id = $1 LIMIT 1;

-- name: SetDebt :one
INSERT INTO debt (
    amount, user_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: UpdateDebt :one
UPDATE debt
SET amount = $1, last_updated = now()
WHERE user_id = $2
RETURNING *;

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
    lower($1)
) RETURNING *;

-- name: NumberOfPlayers :one
SELECT COUNT(*) FROM player;

-- name: GetAllDebts :many
SELECT * FROM player p
JOIN debt d ON p.id = d.user_id
ORDER BY d.amount DESC, upper(p.name);

-- name: GetIdOfPlayer :one
SELECT id FROM player
WHERE name = lower($1) LIMIT 1;

-- name: PutBotSetup :one
INSERT INTO bot_setup (
    channel_id, message_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetBotSetup :one
SELECT * FROM bot_setup
WHERE created_at = (SELECT MAX(created_at) FROM bot_setup) LIMIT 1;