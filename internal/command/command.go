package command

import (
	"context"
	"fmt"
	"slash10k/internal/db"
	"slash10k/internal/utils"
	sqlc "slash10k/sql/gen"
	_ "time/tzdata"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
)

const (
	DeleteReason = "DebtsUpdated"
	TorfstackId  = "263352209654153236"
)

var (
	torfstackId discord.UserID
)

func Setup() error {
	userId, err := discord.ParseSnowflake(TorfstackId)
	if err != nil {
		log.Error().Msgf("cannot parse torfstack id: %s", err)
		return fmt.Errorf("cannot parse torfstack id: %w", err)
	}
	torfstackId = discord.UserID(userId)
	return nil
}

func updateDebtsMessage(ctx context.Context, s *state.State, d db.Database, c DebtClient) {
	debts, err := c.GetAllDebts(ctx)
	if err != nil {
		log.Error().Msgf("cannot get debts: %s", err)
		return
	}

	conn, err := d.Connect(ctx, utils.DefaultConfig().ConnectionString)
	if err != nil {
		log.Error().Msgf("cannot get db connection: %s", err)
	}
	defer func(conn db.Connection, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, ctx)
	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		log.Error().Msgf("cannot start transaction: %s", err)
		return
	}
	defer func(tx db.Transaction, ctx context.Context) {
		_ = tx.Commit(ctx)
	}(tx, ctx)
	queries := tx.Queries()
	botSetup, err := queries.GetBotSetup(ctx)
	if err != nil {
		log.Error().Msgf("cannot get bot setup: %s", err)
		return
	}
	channelId, messageId := botSetupToDiscordTypes(botSetup)

	if channelId == discord.NullChannelID {
		log.Error().Msg("channel id is null")
		return
	}

	if messageId != discord.NullMessageID {
		// Delete old message
		err = s.DeleteMessage(channelId, messageId, DeleteReason)
		if err != nil {
			log.Error().Msgf("cannot delete message: %s", err)
			return
		}
	}

	// Send new message
	m, err := s.SendMessage(channelId, "", *transformDebtsToEmbed(debts))
	if err != nil {
		log.Error().Msgf("cannot edit message: %s", err)
		return
	}

	_, err = queries.PutBotSetup(
		ctx, sqlc.PutBotSetupParams{
			ChannelID: channelId.String(),
			MessageID: m.ID.String(),
		},
	)
}

func botSetupToDiscordTypes(botSetup sqlc.BotSetup) (discord.ChannelID, discord.MessageID) {
	channelId, err := discord.ParseSnowflake(botSetup.ChannelID)
	if err != nil {
		log.Error().Msgf("cannot parse channel id: %s", err)
	}
	messageId, err := discord.ParseSnowflake(botSetup.MessageID)
	if err != nil {
		log.Error().Msgf("cannot parse message id: %s", err)
	}
	return discord.ChannelID(channelId), discord.MessageID(messageId)
}
