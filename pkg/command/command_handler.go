package command

import (
	"context"
	"errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"slash10k/pkg/domain"
)

func RegisterDiscordHandlers(s *state.State, service domain.Service, lookup domain.MessageLookup) {
	s.AddHandler(
		func(event *gateway.MessageReactionAddEvent) {
			ctx := context.Background()
			isRegistrationMessage := lookup.IsRegistrationMessage(event.MessageID.String())
			if isRegistrationMessage {
				if event.Emoji.Name != "💰" {
					return
				}
				log.Info().Msgf("reaction %s added on registration message", event.Emoji.Name)
				err := service.AddPlayer(
					ctx,
					event.UserID.String(),
					event.Member.User.Username,
					event.GuildID.String(),
					event.Member.User.DisplayName,
				)
				if err != nil && !errors.Is(err, domain.ErrPlayerAlreadyExists) {
					log.Error().Msgf("could not add player: %s", err)
					return
				} else if errors.Is(err, domain.ErrPlayerAlreadyExists) {
					log.Warn().Msgf("could not add player: %s", err)
					return
				}
				updateDebtsMessage(ctx, s, service, event.GuildID.String())
			}
		},
	)
	s.AddHandler(
		func(event *gateway.MessageReactionRemoveEvent) {
			ctx := context.Background()
			isRegistrationMessage := lookup.IsRegistrationMessage(event.MessageID.String())
			if isRegistrationMessage {
				if event.Emoji.Name != "💰" {
					return
				}
				log.Info().Msgf("reaction %s removed on registration message", event.Emoji.Name)
				err := service.DeletePlayer(
					ctx,
					event.UserID.String(),
					event.GuildID.String(),
				)
				if err != nil && !errors.Is(err, domain.ErrPlayerDoesNotExist) {
					log.Error().Msgf("could not delete player: %s", err)
					return
				} else if errors.Is(err, domain.ErrPlayerDoesNotExist) {
					log.Warn().Msgf("could not delete player: %s", err)
					return
				}
				updateDebtsMessage(ctx, s, service, event.GuildID.String())
			}
		},
	)
	s.AddHandler(
		func(event *gateway.InteractionCreateEvent) {
			ctx := context.Background()
			switch data := event.Data.(type) {
			case *discord.StringSelectInteraction:
				if data.CustomID == ComponentIdSelectPlayer {
					log.Info().Msgf("select player interaction")
					if len(data.Values) != 1 {
						log.Error().Msgf("invalid number of players selected: %v", len(data.Values))
						return
					}
					player := data.Values[0]
					err := service.AddDebt(ctx, player, event.GuildID.String(), 10000)
					if err != nil {
						log.Error().Msgf("could not add debt: %s", err)
						return
					}
					updateDebtsMessage(ctx, s, service, event.GuildID.String())
					err = s.RespondInteraction(
						event.ID, event.Token, api.InteractionResponse{
							Type: api.UpdateMessage,
							Data: nil,
						},
					)
				}
			default:
				return
			}
		},
	)
	s.AddHandler(
		func(event *gateway.InteractionCreateEvent) {
			ctx := context.Background()
			switch data := event.Data.(type) {
			case *discord.ButtonInteraction:
				if data.CustomID == ComponentIdPaid {
					log.Info().Msgf("paid button interaction")
					err := service.ResetDebt(ctx, event.SenderID().String(), event.GuildID.String())
					if err != nil {
						log.Error().Msgf("could not reset debt: %s", err)
						return
					}
					updateDebtsMessage(ctx, s, service, event.GuildID.String())
					err = s.RespondInteraction(
						event.ID, event.Token, api.InteractionResponse{
							Type: api.UpdateMessage,
							Data: nil,
						},
					)
				}
			default:
				return
			}
		},
	)
}
