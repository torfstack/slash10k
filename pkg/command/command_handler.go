package command

import (
	"context"
	"errors"
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
				log.Info().Msg("reaction added on registration message")
				err := service.AddPlayer(
					ctx,
					event.UserID.String(),
					event.Member.User.Username,
					event.Member.User.DisplayName,
					event.GuildID.String(),
				)
				if err != nil && !errors.Is(err, domain.ErrPlayerAlreadyExists) {
					log.Error().Msgf("could not add player: %s", err)
					return
				} else if errors.Is(err, domain.ErrPlayerAlreadyExists) {
					log.Warn().Msgf("could not add player: %s", err)
					return
				}
				updateDebtsMessage(ctx, s, d, c, event.GuildID.String())
			}
		},
	)
	s.AddHandler(
		func(g *gateway.MessageReactionRemoveEvent) {
			ctx := context.Background()
			isRegistrationMessage := lookup.IsRegistrationMessage(g.MessageID.String())
			if isRegistrationMessage {
				log.Info().Msg("reaction removed on registration message")
				err := service.DeletePlayer(
					ctx,
					g.UserID.String(),
					g.GuildID.String(),
				)
				if err != nil && !errors.Is(err, domain.ErrPlayerDoesNotExist) {
					log.Error().Msgf("could not delete player: %s", err)
					return
				} else if errors.Is(err, domain.ErrPlayerDoesNotExist) {
					log.Warn().Msgf("could not delete player: %s", err)
					return
				}
				updateDebtsMessage(ctx, s, d, c, g.GuildID.String())
			}
		},
	)
}
