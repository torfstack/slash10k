-- name: PutBotSetup :one
INSERT INTO bot_setup (guild_id, channel_id, debts_message_id, registration_message_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteBotSetup :exec
DELETE FROM bot_setup
WHERE guild_id = $1;

-- name: GetBotSetup :one
SELECT *
FROM bot_setup
WHERE bot_setup.guild_id = $1 LIMIT 1;

-- name: DoesBotSetupExist :one
SELECT EXISTS(SELECT 1 FROM bot_setup WHERE guild_id = $1);

-- name: GetAllBotSetups :many
SELECT *
FROM bot_setup;