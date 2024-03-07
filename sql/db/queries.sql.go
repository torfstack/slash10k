// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addDebtJournalEntry = `-- name: AddDebtJournalEntry :one
INSERT INTO debt_journal (
    amount, description, user_id
) VALUES (
    $1, $2, $3
) RETURNING id, amount, description, date, user_id
`

type AddDebtJournalEntryParams struct {
	Amount      int64
	Description string
	UserID      pgtype.Int4
}

func (q *Queries) AddDebtJournalEntry(ctx context.Context, arg AddDebtJournalEntryParams) (DebtJournal, error) {
	row := q.db.QueryRow(ctx, addDebtJournalEntry, arg.Amount, arg.Description, arg.UserID)
	var i DebtJournal
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Description,
		&i.Date,
		&i.UserID,
	)
	return i, err
}

const addPlayer = `-- name: AddPlayer :one
INSERT INTO player (
    name
) VALUES (
    $1
) RETURNING id, name
`

func (q *Queries) AddPlayer(ctx context.Context, name string) (Player, error) {
	row := q.db.QueryRow(ctx, addPlayer, name)
	var i Player
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const allPlayerDebts = `-- name: AllPlayerDebts :many
SELECT p.id, name, d.id, amount, last_updated, user_id FROM player p
JOIN debt d ON p.id = d.user_id
ORDER BY d.amount DESC, upper(p.name)
`

type AllPlayerDebtsRow struct {
	ID          int32
	Name        string
	ID_2        int32
	Amount      int64
	LastUpdated pgtype.Timestamp
	UserID      pgtype.Int4
}

func (q *Queries) AllPlayerDebts(ctx context.Context) ([]AllPlayerDebtsRow, error) {
	rows, err := q.db.Query(ctx, allPlayerDebts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AllPlayerDebtsRow
	for rows.Next() {
		var i AllPlayerDebtsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ID_2,
			&i.Amount,
			&i.LastUpdated,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDebt = `-- name: GetDebt :one
SELECT id, amount, last_updated, user_id FROM debt
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetDebt(ctx context.Context, userID pgtype.Int4) (Debt, error) {
	row := q.db.QueryRow(ctx, getDebt, userID)
	var i Debt
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.LastUpdated,
		&i.UserID,
	)
	return i, err
}

const getIdOfPlayer = `-- name: GetIdOfPlayer :one
SELECT id FROM player
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetIdOfPlayer(ctx context.Context, name string) (int32, error) {
	row := q.db.QueryRow(ctx, getIdOfPlayer, name)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const numberOfPlayers = `-- name: NumberOfPlayers :one
SELECT COUNT(*) FROM player
`

func (q *Queries) NumberOfPlayers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, numberOfPlayers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const setDebt = `-- name: SetDebt :one
INSERT INTO debt (
    amount, user_id
) VALUES (
    $1, $2
) RETURNING id, amount, last_updated, user_id
`

type SetDebtParams struct {
	Amount int64
	UserID pgtype.Int4
}

func (q *Queries) SetDebt(ctx context.Context, arg SetDebtParams) (Debt, error) {
	row := q.db.QueryRow(ctx, setDebt, arg.Amount, arg.UserID)
	var i Debt
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.LastUpdated,
		&i.UserID,
	)
	return i, err
}

const updateDebt = `-- name: UpdateDebt :one
UPDATE debt
SET amount = $1, last_updated = now()
WHERE user_id = $2
RETURNING id, amount, last_updated, user_id
`

type UpdateDebtParams struct {
	Amount int64
	UserID pgtype.Int4
}

func (q *Queries) UpdateDebt(ctx context.Context, arg UpdateDebtParams) (Debt, error) {
	row := q.db.QueryRow(ctx, updateDebt, arg.Amount, arg.UserID)
	var i Debt
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.LastUpdated,
		&i.UserID,
	)
	return i, err
}
