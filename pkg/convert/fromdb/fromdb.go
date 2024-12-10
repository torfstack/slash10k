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

func FromPlayerWithDebt(playerWithDebt sqlc.GetPlayerRow) models.Player {
	player := FromPlayerWithoutDebt(playerWithDebt.Player)
	player.Debt = FromDebt(playerWithDebt.Debt)
	return player
}

func FromAllPlayers(allPlayers []sqlc.GetAllPlayersRow) []models.Player {
	players := make([]models.Player, len(allPlayers))
	for i, player := range allPlayers {
		p := FromPlayerWithoutDebt(player.Player)
		p.Debt = FromDebt(player.Debt)
		players[i] = p
	}
	return players
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
