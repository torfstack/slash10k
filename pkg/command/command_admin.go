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
	"slash10k/pkg/domain"
)

func SetChannel(state *state.State, service domain.Service, lookup domain.MessageLookup) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		guildId := data.Event.GuildID
		log.Info().Msgf("setup channel called for guild %s", guildId)

		if data.Event.SenderID() != config.TorfstackUserId() {
			log.Error().Msgf("cannot set channel: not torfstack, got %v", data.Event.SenderID())
			return ephemeralMessage("You are not allowed to set the channel, ask Torfstack!")
		}
		log.Debug().Msg("called by correct user 'torfstack'")

		options := data.Options
		var err error
		cId, err := options.Find("channel_id").SnowflakeValue()
		if err != nil {
			log.Error().Msgf("cannot get channel_id: %s", err)
			return ephemeralMessage("Could not set channel")
		}

		channelId := discord.ChannelID(cId)

		alreadySetup, err := isAlreadySetup(ctx, service, guildId.String())
		if err != nil {
			log.Error().Msgf("cannot check if already setup: %s", err)
			return ephemeralMessage("Could not check if already setup")
		}
		if alreadySetup {
			log.Debug().Msgf("already setup for guild %s, deleting messages and current setup", guildId)
			err = deleteMessagesAndCurrentSetup(ctx, state, service, guildId.String())
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
		log.Debug().Msgf("registration message sent with id: %s", registrationMessageId)

		debtsMessage, err := sendDebtsMessage(ctx, state, service, guildId.String(), channelId)
		if err != nil {
			log.Error().Msgf("cannot send debts message: %s", err)
			return ephemeralMessage("Could not send debts message")
		}
		debtsMessageId := debtsMessage.ID
		log.Debug().Msgf("debts message sent with id: %s", debtsMessageId)

		err = service.SetBotSetup(
			ctx,
			guildId.String(),
			channelId.String(),
			registrationMessageId.String(),
			debtsMessageId.String(),
		)
		if err != nil {
			log.Error().Msgf("cannot put bot setup: %s", err)
			return ephemeralMessage("Could not put bot setup")
		}

		botSetup, err := service.GetBotSetup(ctx, guildId.String())
		if err != nil {
			log.Error().Msgf("cannot get bot setup: %s", err)
			return ephemeralMessage("Could not get bot setup")
		}
		lookup.AddSetup(*botSetup)

		return ephemeralMessage("Channel set successfully")
	}
}

func deleteMessagesAndCurrentSetup(ctx context.Context, s *state.State, service domain.Service, guildId string) error {
	botSetup, err := service.GetBotSetup(ctx, guildId)
	if err != nil {
		return errors.New("could not get bot setup")
	}
	channelId, err := discord.ParseSnowflake(botSetup.ChannelId)
	if err != nil {
		return errors.New("could not parse channel id")
	}
	debtsMessageId, err := discord.ParseSnowflake(botSetup.DebtsMessageId)
	if err != nil {
		return errors.New("could not parse debts message id")
	}
	registrationMessageId, err := discord.ParseSnowflake(botSetup.RegistrationMessageId)
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
	return service.DeleteBotSetup(ctx, guildId)
}

func isAlreadySetup(
	ctx context.Context,
	service domain.Service,
	guildId string,
) (bool, error) {
	_, err := service.GetBotSetup(ctx, guildId)
	if err != nil && !errors.Is(err, domain.ErrBotSetupDoesNotExist) {
		return false, errors.New("could not get bot setup")
	} else if errors.Is(err, domain.ErrBotSetupDoesNotExist) {
		return false, nil
	}
	return true, nil
}

func sendRegistrationMessage(
	s *state.State,
	channelId discord.ChannelID,
) (*discord.Message, error) {
	m, err := s.SendMessage(channelId, ":moneybag: react to join!")
	if err != nil {
		return nil, errors.New("could not send message")
	}
	return m, nil
}

func sendDebtsMessage(
	ctx context.Context,
	s *state.State,
	service domain.Service,
	guildId string,
	channelId discord.ChannelID,
) (*discord.Message, error) {
	allPlayers, err := service.GetAllPlayers(ctx, guildId)
	if err != nil {
		return nil, errors.New("could not get all players")
	}

	m, err := s.SendMessageComplex(channelId, debtsForSendMessage(allPlayers))
	if err != nil {
		return nil, errors.New("could not send message")
	}
	return m, nil
}
