package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"slash10k/internal/db"
	"slash10k/internal/utils"
	sqlc "slash10k/sql/gen"
)

func AddPlayer(s *state.State, d db.Database, c AdminClient) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		if data.Event.SenderID() != torfstackId {
			log.Error().Msgf("cannot add player: not torfstack, got %v", data.Event.SenderID())
			return ephemeralMessage("You are not allowed to add a player, ask Torfstack!")
		}
		options := data.Options
		name := options.Find("name").String()

		err := c.AddPlayer(ctx, name)
		if err != nil {
			log.Error().Msgf("could not add player: %s", err)
			return ephemeralMessage("Could not add player")
		}

		updateDebtsMessage(ctx, s, d, c)
		return ephemeralMessage(fmt.Sprintf("Added player %v", name))
	}
}

func DeletePlayer(s *state.State, d db.Database, c AdminClient) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		if data.Event.SenderID() != torfstackId {
			log.Error().Msgf("cannot delete player: not torfstack, got %v", data.Event.SenderID())
			return ephemeralMessage("You are not allowed to delete a player, ask Torfstack!")
		}
		options := data.Options
		name := options.Find("name").String()
		updateDebtsMessage(ctx, s, d, c)
		return ephemeralMessage(fmt.Sprintf("Deleted player %v", name))
	}
}

func SetChannel(s *state.State, d db.Database, c DebtClient) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		if data.Event.SenderID() != torfstackId {
			log.Error().Msgf("cannot set channel: not torfstack, got %v", data.Event.SenderID())
			return ephemeralMessage("You are not allowed to set the channel, ask Torfstack!")
		}
		options := data.Options
		var err error
		cId, err := options.Find("channel_id").SnowflakeValue()
		if err != nil {
			log.Error().Msgf("cannot get channel_id: %s", err)
			return ephemeralMessage("Could not set channel")
		}

		channelId := discord.ChannelID(cId)
		debts, err := c.GetAllDebts(ctx)
		if err != nil {
			log.Error().Msgf("cannot get debts: %s", err)
			return ephemeralMessage("Could not get debts")
		}
		m, err := s.SendMessage(channelId, "", *transformDebtsToEmbed(debts))
		if err != nil {
			log.Error().Msgf("cannot send message: %s", err)
			return ephemeralMessage("Could not send message")
		}
		messageId := m.ID
		conn, err := d.Connect(ctx, utils.DefaultConfig().ConnectionString)
		if err != nil {
			log.Error().Msgf("cannot get db connection: %s", err)
			return ephemeralMessage("Could not get db connection")
		}
		defer func(conn db.Connection, ctx context.Context) {
			_ = conn.Close(ctx)
		}(conn, ctx)
		_, err = conn.Queries().PutBotSetup(
			ctx, sqlc.PutBotSetupParams{
				ChannelID: channelId.String(),
				MessageID: messageId.String(),
			},
		)
		if err != nil {
			log.Error().Msgf("cannot put bot setup: %s", err)
			return ephemeralMessage("Could not put bot setup")
		}
		return ephemeralMessage("Channel set successfully")
	}
}
