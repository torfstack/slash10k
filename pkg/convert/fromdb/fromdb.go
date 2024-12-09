package fromdb

import (
	"slash10k/pkg/models"
	sqlc "slash10k/sql/gen"
)

func FromPlayerWithoutDebt(player sqlc.Player) models.Player {
	return models.Player{
		Id:          player.ID,
		DiscordId:   player.DiscordID,
		DiscordName: player.DiscordName,
		GuildId:     player.GuildID,
		Name:        player.Name,
	}
}

func FromPlayerWithDebt(playerWithDebt []sqlc.GetPlayerRow) models.Player {
	player := FromPlayerWithoutDebt(playerWithDebt[0].Player)
	player.Debt = FromDebt(playerWithDebt[0].Debt)
	player.DebtJournal = make([]models.DebtJournalEntry, len(playerWithDebt))
	for i, row := range playerWithDebt {
		player.DebtJournal[i] = FromDebtJournal(row.DebtJournal)
	}
	return player
}

func FromAllPlayers(allPlayers []sqlc.GetAllPlayersRow) []models.Player {
	players := make(map[string]models.Player, len(allPlayers))
	for _, row := range allPlayers {
		p, ok := players[row.Player.DiscordID]
		if !ok {
			p = FromPlayerWithoutDebt(row.Player)
			p.Debt = FromDebt(row.Debt)
			p.DebtJournal = make([]models.DebtJournalEntry, 10)
			p.DebtJournal[0] = FromDebtJournal(row.DebtJournal)
			players[row.Player.DiscordID] = p
		} else {
			p.DebtJournal = append(p.DebtJournal, FromDebtJournal(row.DebtJournal))
		}
	}
	res := make([]models.Player, len(players))
	for _, player := range players {
		res = append(res, player)
	}
	return res
}

func FromDebt(debt sqlc.Debt) models.Debt {
	return models.Debt{
		Id:          debt.ID,
		Amount:      debt.Amount,
		LastUpdated: debt.LastUpdated.Time.Unix(),
		UserId:      debt.UserID,
	}
}

func FromDebtJournal(debtJournal sqlc.DebtJournal) models.DebtJournalEntry {
	return models.DebtJournalEntry{
		Id:          debtJournal.ID,
		Amount:      debtJournal.Amount,
		Description: debtJournal.Description,
		Date:        debtJournal.Date.Time.Unix(),
		UserId:      debtJournal.UserID,
	}
}

func FromBotSetup(botSetup sqlc.BotSetup) models.BotSetup {
	return models.BotSetup{
		GuildId:               botSetup.GuildID,
		ChannelId:             botSetup.ChannelID,
		RegistrationMessageId: botSetup.RegistrationMessageID,
		DebtsMessageId:        botSetup.DebtsMessageID,
	}
}

func FromBotSetups(botSetups []sqlc.BotSetup) []models.BotSetup {
	botSetupsConverted := make([]models.BotSetup, len(botSetups))
	for i, botSetup := range botSetups {
		botSetupsConverted[i] = FromBotSetup(botSetup)
	}
	return botSetupsConverted
}
