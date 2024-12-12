package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"os"
	"slash10k/pkg/domain"
	"strings"
)

var (
	tokenUuidMap = make(map[string]string)
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
			case *discord.ButtonInteraction:
				switch {
				case data.CustomID == ComponentIdPaid:
					log.Info().Msgf("paid button interaction")
					err := service.ResetDebt(ctx, event.SenderID().String(), event.GuildID.String())
					if err != nil {
						log.Error().Msgf("could not reset debt: %s", err)
						return
					}
					updateDebtsMessage(ctx, s, service, event.GuildID.String())
				case strings.HasPrefix(string(data.CustomID), ComponentIdCancelButton):
					appIdSnowflake, err := discord.ParseSnowflake(os.Getenv("APPLICATION_ID"))
					if err != nil {
						log.Error().Msgf("could not parse application id: %s", err)
						return
					}
					updateDebtsMessage(ctx, s, service, event.GuildID.String())
					_, originalToken := extractPlayerAndToken(string(data.CustomID))
					err = s.DeleteInteractionResponse(discord.AppID(appIdSnowflake), originalToken)
					if err != nil {
						log.Error().Msgf("could not delete interaction response: %s", err)
						return
					}
				case strings.HasPrefix(string(data.CustomID), ComponentIdConfirmButton):
					player, originalToken := extractPlayerAndToken(string(data.CustomID))
					err := service.AddDebt(ctx, player, event.GuildID.String(), 10000)
					if err != nil {
						log.Error().Msgf("could not add debt: %s", err)
						return
					}
					updateDebtsMessage(ctx, s, service, event.GuildID.String())
					appIdSnowflake, err := discord.ParseSnowflake(os.Getenv("APPLICATION_ID"))
					if err != nil {
						log.Error().Msgf("could not parse application id: %s", err)
						return
					}
					err = s.DeleteInteractionResponse(discord.AppID(appIdSnowflake), originalToken)
					if err != nil {
						log.Error().Msgf("could not delete interaction response: %s", err)
						return
					}
				}
				err := s.RespondInteraction(
					event.ID, event.Token, api.InteractionResponse{
						Type: api.DeferredMessageUpdate,
						Data: nil,
					},
				)
				if err != nil {
					log.Error().Msgf("could not respond to interaction: %s", err)
					return
				}
			case *discord.StringSelectInteraction:
				if data.CustomID == ComponentIdSelectPlayer {
					log.Info().Msgf("select player interaction")
					if len(data.Values) != 1 {
						log.Error().Msgf("invalid number of players selected: %v", len(data.Values))
						return
					}
					player, err := service.GetPlayer(ctx, data.Values[0], event.GuildID.String())
					if err != nil {
						log.Error().Msgf("could not get player: %s", err)
						return
					}
					u := uuid.NewString()
					tokenUuidMap[u] = event.Token
					err = s.RespondInteraction(
						event.ID, event.Token, api.InteractionResponse{
							Type: api.MessageInteractionWithSource,
							Data: &api.InteractionResponseData{
								Content:    option.NewNullableString("Do you really want to add 10k to " + player.Name + "?"),
								Components: confirmOrCancelButtonComponents(player.DiscordId, u),
								Flags:      discord.EphemeralMessage,
							},
						},
					)
					if err != nil {
						log.Error().Msgf("could not respond to interaction: %s", err)
						return
					}
				}
			default:
				return
			}
		},
	)
}

const (
	ComponentIdCancelButton     = "CANCEL"
	ComponentLabelCancelButton  = "Cancel"
	ComponentIdConfirmButton    = "CONFIRM"
	ComponentLabelConfirmButton = "Confirm"
)

func confirmOrCancelButtonComponents(player string, token string) *discord.ContainerComponents {
	return &discord.ContainerComponents{
		&discord.ActionRowComponent{
			&discord.ButtonComponent{
				Style:    discord.SecondaryButtonStyle(),
				CustomID: discord.ComponentID(fmt.Sprintf("%s||%s||%s", ComponentIdCancelButton, player, token)),
				Label:    ComponentLabelCancelButton,
			},
			&discord.ButtonComponent{
				Style:    discord.DangerButtonStyle(),
				CustomID: discord.ComponentID(fmt.Sprintf("%s||%s||%s", ComponentIdConfirmButton, player, token)),
				Label:    ComponentLabelConfirmButton,
			},
		},
	}
}

func extractPlayerAndToken(customId string) (string, string) {
	playerAndToken := strings.TrimPrefix(customId, ComponentIdCancelButton+"||")
	playerAndToken = strings.TrimPrefix(playerAndToken, ComponentIdConfirmButton+"||")
	components := strings.Split(playerAndToken, "||")
	return components[0], tokenUuidMap[components[1]]
}
