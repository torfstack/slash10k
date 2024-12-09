package command

import (
	"context"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"slash10k/pkg/db"
	"slash10k/pkg/domain"
	"slash10k/pkg/models"
	sqlc "slash10k/sql/gen"
	_ "time/tzdata"
)

func updateDebtsMessage(ctx context.Context, discord *state.State, service domain.Service, guildId string) {
	debts, err := allDebts(ctx, conn, guildId)
	if err != nil {
		log.Error().Msgf("cannot get debts: %s", err)
		return
	}

	queries := tx.Queries()
	botSetup, err := queries.GetBotSetup(ctx, guildId)
	if err != nil {
		log.Error().Msgf("cannot get bot setup: %s", err)
		return
	}
	channelId, messageId := botSetupToDiscordTypes(botSetup)

	if channelId == discord.NullChannelID || messageId == discord.NullMessageID {
		log.Error().Msg("channel id is null")
		return
	}

	// Edit debts message
	_, err = discord.EditMessage(
		channelId,
		messageId,
		"",
		*transformDebtsToEmbed(&models.AllDebtsResponse{Debts: debts}),
	)
	if err != nil {
		log.Error().Msgf("cannot edit message: %s", err)
		return
	}
}

func botSetupToDiscordTypes(botSetup sqlc.BotSetup) (discord.ChannelID, discord.MessageID) {
	channelId, err := discord.ParseSnowflake(botSetup.ChannelID)
	if err != nil {
		log.Error().Msgf("cannot parse channel id: %s", err)
	}
	messageId, err := discord.ParseSnowflake(botSetup.DebtsMessageID)
	if err != nil {
		log.Error().Msgf("cannot parse message id: %s", err)
	}
	return discord.ChannelID(channelId), discord.MessageID(messageId)
}
