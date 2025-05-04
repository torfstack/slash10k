-- name: SetDebt :exec
INSERT INTO debt (amount, user_id)
VALUES ($1, $2)
ON CONFLICT (user_id)
    DO UPDATE SET amount = $1, last_updated = now()
WHERE debt.user_id = $2;

