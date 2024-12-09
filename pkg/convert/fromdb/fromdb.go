package fromdb

import (
	"slash10k/pkg/domain"
	sqlc "slash10k/sql/gen"
)

func FromPlayerWithoutDebt(player sqlc.Player) domain.Player {
	return domain.Player{
		Id:          player.ID,
		DiscordId:   player.DiscordID,
		DiscordName: player.DiscordName,
		GuildId:     player.GuildID,
		Name:        player.Name,
	}
}

func FromPlayerWithDebt(playerWithDebt []sqlc.GetPlayerRow) domain.Player {
	player := FromPlayerWithoutDebt(playerWithDebt[0].Player)
	player.Debt = FromDebt(playerWithDebt[0].Debt)
	player.DebtJournal = make([]domain.DebtJournalEntry, len(playerWithDebt))
	for i, row := range playerWithDebt {
		player.DebtJournal[i] = FromDebtJournal(row.DebtJournal)
	}
	return player
}

func FromAllPlayers(allPlayers []sqlc.GetAllPlayersRow) []domain.Player {
	players := make(map[string]domain.Player, len(allPlayers))
	for _, row := range allPlayers {
		p, ok := players[row.Player.DiscordID]
		if !ok {
			p = FromPlayerWithoutDebt(row.Player)
			p.Debt = FromDebt(row.Debt)
			p.DebtJournal = make([]domain.DebtJournalEntry, 10)
			p.DebtJournal[0] = FromDebtJournal(row.DebtJournal)
			players[row.Player.DiscordID] = p
		} else {
			p.DebtJournal = append(p.DebtJournal, FromDebtJournal(row.DebtJournal))
		}
	}
	res := make([]domain.Player, len(players))
	for _, player := range players {
		res = append(res, player)
	}
	return res
}

func FromDebt(debt sqlc.Debt) domain.Debt {
	return domain.Debt{
		Id:          debt.ID,
		Amount:      debt.Amount,
		LastUpdated: debt.LastUpdated.Time.Unix(),
		UserId:      debt.UserID,
	}
}

func FromDebtJournal(debtJournal sqlc.DebtJournal) domain.DebtJournalEntry {
	return domain.DebtJournalEntry{
		Id:          debtJournal.ID,
		Amount:      debtJournal.Amount,
		Description: debtJournal.Description,
		Date:        debtJournal.Date.Time.Unix(),
		UserId:      debtJournal.UserID,
	}
}

func FromBotSetup(botSetup sqlc.BotSetup) domain.BotSetup {
	return domain.BotSetup{
		GuildId:               botSetup.GuildID,
		ChannelId:             botSetup.ChannelID,
		RegistrationMessageId: botSetup.RegistrationMessageID,
		DebtsMessageId:        botSetup.DebtsMessageID,
	}
}

func FromBotSetups(botSetups []sqlc.BotSetup) []domain.BotSetup {
	botSetupsConverted := make([]domain.BotSetup, len(botSetups))
	for i, botSetup := range botSetups {
		botSetupsConverted[i] = FromBotSetup(botSetup)
	}
	return botSetupsConverted
}
