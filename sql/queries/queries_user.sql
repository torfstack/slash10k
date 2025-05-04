-- name: DoesPlayerExist :one
SELECT EXISTS(SELECT 1 FROM player WHERE discord_id = $1 AND guild_id = $2);

-- name: AddPlayer :one
INSERT INTO player (discord_id, discord_name, guild_id, name)
VALUES ($1, $2, $3, $4 )
RETURNING *;

-- name: GetPlayer :one
SELECT sqlc.embed(player), sqlc.embed(debt)
FROM player
JOIN debt ON player.id = debt.user_id
WHERE player.discord_id = $1 AND player.guild_id = $2
LIMIT 1;

-- name: GetAllPlayers :many
SELECT sqlc.embed(player), sqlc.embed(debt)
FROM player
JOIN debt ON player.id = debt.user_id
WHERE player.guild_id = $1;

-- name: DeletePlayer :exec
DELETE FROM player
WHERE id = $1;

-- name: NumberOfPlayers :one
SELECT COUNT(discord_id)
FROM player;

-- name: GetIdOfPlayer :one
SELECT id FROM player
WHERE discord_id = $1 AND guild_id = $2
LIMIT 1;
