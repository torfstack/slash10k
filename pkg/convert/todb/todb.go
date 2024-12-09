package todb

import (
	"github.com/jackc/pgx/v5/pgtype"
	"slash10k/pkg/domain"
	sqlc "slash10k/sql/gen"
	"time"
)

func ToPlayer(player domain.Player) sqlc.Player {
	return sqlc.Player{
		ID:          player.Id,
		DiscordID:   player.DiscordId,
		DiscordName: player.DiscordName,
		GuildID:     player.GuildId,
		Name:        player.Name,
	}
}

func ToDebt(debt domain.Debt) sqlc.Debt {
	return sqlc.Debt{
		ID:     debt.Id,
		Amount: debt.Amount,
		LastUpdated: pgtype.Timestamp{
			Time:  time.Unix(debt.LastUpdated, 0),
			Valid: true,
		},
		UserID: debt.UserId,
	}
}

func ToDebtJournal(debtJournal domain.DebtJournalEntry) sqlc.DebtJournal {
	return sqlc.DebtJournal{
		ID:          debtJournal.Id,
		Amount:      debtJournal.Amount,
		Description: debtJournal.Description,
		Date: pgtype.Timestamp{
			Time:  time.Unix(debtJournal.Date, 0),
			Valid: true,
		},
		UserID: debtJournal.UserId,
	}
}
