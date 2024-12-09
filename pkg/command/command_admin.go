package command

import (
	"context"
	"errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"slash10k/pkg/config"
	"slash10k/pkg/convert/fromdb"
	"slash10k/pkg/db"
	"slash10k/pkg/domain"
	"slash10k/pkg/models"
	sqlc "slash10k/sql/gen"
)

func SetChannel(state *state.State, service domain.Service, lookup domain.MessageLookup) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		if data.Event.SenderID() != config.TorfstackUserId() {
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

		guildId := data.Event.GuildID
		channelId := discord.ChannelID(cId)

		alreadySetup, err := isAlreadySetup(ctx, conn, guildId.String())
		if err != nil {
			log.Error().Msgf("cannot check if already setup: %s", err)
			return ephemeralMessage("Could not check if already setup")
		}
		if alreadySetup {
			err = deleteMessagesAndCurrentSetup(ctx, state, conn, guildId.String())
			if err != nil {
				log.Error().Msgf("cannot delete messages and current setup: %s", err)
				return ephemeralMessage("Could not delete messages and current setup")
			}
		}

		registrationMessage, err := sendRegistrationMessage(state, channelId)
		if err != nil {
			log.Error().Msgf("cannot send registration message: %s", err)
			return ephemeralMessage("Could not send registration message")
		}
		registrationMessageId := registrationMessage.ID

		debtsMessage, err := sendDebtsMessage(ctx, state, conn, guildId.String(), channelId)
		if err != nil {
			log.Error().Msgf("cannot send debts message: %s", err)
			return ephemeralMessage("Could not send debts message")
		}
		debtsMessageId := debtsMessage.ID

		botSetup, err := conn.Queries().PutBotSetup(
			ctx, sqlc.PutBotSetupParams{
				ChannelID:             channelId.String(),
				DebtsMessageID:        debtsMessageId.String(),
				GuildID:               guildId.String(),
				RegistrationMessageID: registrationMessageId.String(),
			},
		)
		if err != nil {
			log.Error().Msgf("cannot put bot setup: %s", err)
			return ephemeralMessage("Could not put bot setup")
		}

		lookup.AddSetup(fromdb.FromBotSetup(botSetup))

		return ephemeralMessage("Channel set successfully")
	}
}

func deleteMessagesAndCurrentSetup(ctx context.Context, s *state.State, conn db.Connection, guildId string) error {
	botSetup, err := conn.Queries().GetBotSetup(ctx, guildId)
	if err != nil {
		return errors.New("could not get bot setup")
	}
	channelId, err := discord.ParseSnowflake(botSetup.ChannelID)
	if err != nil {
		return errors.New("could not parse channel id")
	}
	debtsMessageId, err := discord.ParseSnowflake(botSetup.DebtsMessageID)
	if err != nil {
		return errors.New("could not parse debts message id")
	}
	registrationMessageId, err := discord.ParseSnowflake(botSetup.RegistrationMessageID)
	if err != nil {
		return errors.New("could not parse registration message id")
	}
	err = s.MessageRemove(
		discord.ChannelID(channelId),
		discord.MessageID(debtsMessageId),
	)
	if err != nil {
		return errors.New("could not delete debts message")
	}
	err = s.MessageRemove(
		discord.ChannelID(channelId),
		discord.MessageID(registrationMessageId),
	)
	if err != nil {
		return errors.New("could not delete registration message")
	}
	return conn.Queries().DeleteBotSetup(ctx, guildId)
}

func isAlreadySetup(
	ctx context.Context,
	conn db.Connection,
	guildId string,
) (bool, error) {
	botSetup, err := conn.Queries().DoesBotSetupExist(ctx, guildId)
	if err != nil {
		return false, errors.New("could not get bot setup")
	}
	return botSetup, nil
}

func sendRegistrationMessage(
	s *state.State,
	channelId discord.ChannelID,
) (*discord.Message, error) {
	m, err := s.SendMessage(channelId, "Hier mal Emoji drauf!")
	if err != nil {
		return nil, errors.New("could not send message")
	}
	return m, nil
}

func sendDebtsMessage(
	ctx context.Context,
	s *state.State,
	conn db.Connection,
	guildId string,
	channelId discord.ChannelID,
) (*discord.Message, error) {
	debts, err := allDebts(ctx, conn, guildId)
	if err != nil {
		return nil, errors.New("could not get debts")
	}

	m, err := s.SendMessage(channelId, "", *transformDebtsToEmbed(&models.AllDebtsResponse{Debts: debts}))
	if err != nil {
		return nil, errors.New("could not send message")
	}
	return m, nil
}
