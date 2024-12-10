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

-- name: DeleteJournalEntry :exec
DELETE FROM debt_journal
WHERE id = $1;

-- name: GetJournalEntries :many
SELECT * FROM debt_journal
WHERE user_id = $1;

-- name: DoesPlayerExist :one
SELECT EXISTS(SELECT 1 FROM player WHERE discord_id = $1 AND guild_id = $2);

-- name: AddPlayer :one
INSERT INTO player (
    discord_id, discord_name, guild_id, name
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPlayer :many
SELECT sqlc.embed(player), sqlc.embed(debt) FROM player
JOIN debt ON player.id = debt.user_id
WHERE player.discord_id = $1 AND player.guild_id = $2;

-- name: GetAllPlayers :many
SELECT sqlc.embed(player), sqlc.embed(debt) FROM player
JOIN debt ON player.id = debt.user_id
WHERE guild_id = $1;

-- name: DeletePlayer :exec
DELETE FROM player
WHERE id = $1;

-- name: NumberOfPlayers :one
SELECT COUNT(discord_id) FROM player;

-- name: GetIdOfPlayer :one
SELECT id FROM player
WHERE discord_id = $1 AND guild_id = $2 LIMIT 1;

-- name: PutBotSetup :one
INSERT INTO bot_setup (
    guild_id, channel_id, debts_message_id, registration_message_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: DeleteBotSetup :exec
DELETE FROM bot_setup
WHERE guild_id = $1;

-- name: GetBotSetup :one
SELECT * FROM bot_setup
WHERE created_at = (SELECT MAX(created_at) FROM bot_setup) AND bot_setup.guild_id = $1 LIMIT 1;

-- name: DoesBotSetupExist :one
SELECT EXISTS(SELECT 1 FROM bot_setup WHERE guild_id = $1);

-- name: GetAllBotSetups :many
SELECT * FROM bot_setup;