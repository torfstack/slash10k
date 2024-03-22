-- name: GetDebt :one
SELECT * FROM debt
WHERE user_id = $1 LIMIT 1;

-- name: SetDebt :exec
INSERT INTO debt (amount, user_id)
VALUES ($1, $2)
ON CONFLICT (user_id)
DO UPDATE SET amount = $1, last_updated = now()
WHERE debt.user_id = $2;

-- name: AddJournalEntry :one
INSERT INTO debt_journal (
    amount, description, user_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateJournalEntry :one
UPDATE debt_journal
SET amount = $1, description = $2
WHERE id = $3
RETURNING *;

-- name: GetJournalEntries :many
SELECT * FROM debt_journal
WHERE user_id = $1;

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