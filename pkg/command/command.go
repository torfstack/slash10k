package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/rs/zerolog/log"
	"os"
	"slash10k/pkg/domain"
	"slash10k/pkg/models"
	"slices"
	"sort"
	"strings"
	_ "time/tzdata"
)

const (
	ComponentIdSelectPlayer          = "SELECT_PLAYER"
	ComponentPlaceholderSelectPlayer = "Select a player"
	ComponentIdPaid                  = "PAID"
	ComponentLabelPaid               = "I paid!"
)

func updateDebtsMessage(ctx context.Context, state *state.State, service domain.Service, guildId string) {
	allPlayers, err := service.GetAllPlayers(ctx, guildId)
	if err != nil {
		log.Error().Msgf("cannot get all players: %s", err)
		return
	}
	log.Debug().Msgf("retrieved all players (length %v) for guild %s", len(allPlayers), guildId)

	botSetup, err := service.GetBotSetup(ctx, guildId)
	if err != nil {
		log.Error().Msgf("cannot get bot setup: %s", err)
		return
	}
	log.Debug().Msgf("retrieved bot-setup for guild %s", guildId)
	channelId, messageId := botSetupToDiscordTypes(*botSetup)

	if channelId == discord.NullChannelID || messageId == discord.NullMessageID {
		log.Error().Msg("channel id or message id is null")
		return
	}

	_, err = state.EditMessageComplex(
		channelId,
		messageId,
		debtsForEditMessage(allPlayers),
	)
	if err != nil {
		log.Error().Msgf("cannot edit message: %s", err)
		return
	}
	log.Debug().Msgf("edited debt message for guild %s", guildId)
}

func botSetupToDiscordTypes(botSetup models.BotSetup) (discord.ChannelID, discord.MessageID) {
	channelId, err := discord.ParseSnowflake(botSetup.ChannelId)
	if err != nil {
		log.Error().Msgf("cannot parse channel id: %s", err)
	}
	messageId, err := discord.ParseSnowflake(botSetup.DebtsMessageId)
	if err != nil {
		log.Error().Msgf("cannot parse message id: %s", err)
	}
	return discord.ChannelID(channelId), discord.MessageID(messageId)
}

func debtsForSendMessage(allPlayers []models.Player) api.SendMessageData {
	return api.SendMessageData{
		Content:    "",
		Embeds:     []discord.Embed{transformDebtsToEmbed(allPlayers)},
		Components: buttonComponents(allPlayers),
	}
}

func debtsForEditMessage(allPlayers []models.Player) api.EditMessageData {
	buttons := buttonComponents(allPlayers)
	return api.EditMessageData{
		Content:    option.NewNullableString(""),
		Embeds:     &[]discord.Embed{transformDebtsToEmbed(allPlayers)},
		Components: &buttons,
	}
}

func buttonComponents(allPlayers []models.Player) discord.ContainerComponents {
	if len(allPlayers) == 0 {
		return make(discord.ContainerComponents, 0)
	}
	playerNames := make([]discord.SelectOption, len(allPlayers))
	for i, p := range allPlayers {
		playerNames[i] = discord.SelectOption{
			Label: p.Name,
			Value: p.DiscordId,
		}
	}
	sort.Slice(
		playerNames, func(i, j int) bool {
			return strings.Compare(playerNames[i].Label, playerNames[j].Label) < 0
		},
	)
	return discord.ContainerComponents{
		&discord.ActionRowComponent{
			&discord.StringSelectComponent{
				Options:     playerNames,
				CustomID:    ComponentIdSelectPlayer,
				Placeholder: ComponentPlaceholderSelectPlayer,
			},
		},
		&discord.ActionRowComponent{
			&discord.ButtonComponent{
				Style:    discord.PrimaryButtonStyle(),
				CustomID: ComponentIdPaid,
				Label:    ComponentLabelPaid,
			},
		},
	}
}

func transformDebtsToEmbed(players []models.Player) discord.Embed {
	embed := defaultEmbed()

	if len(players) > 0 {
		maxLength := len(
			slices.MaxFunc(
				players, func(p1, p2 models.Player) int {
					return len(p1.DiscordName) - len(p2.DiscordName)
				},
			).DiscordName,
		)
		debtString := strings.Builder{}
		debtString.WriteString("```")
		for _, p := range players {
			debtString.WriteString(fmt.Sprintf("%-*s %v\n", maxLength, p.Name, p.Debt.Amount))
		}
		debtString.WriteString("```")
		embed.Fields = []discord.EmbedField{
			{
				Name:   "Spieler",
				Value:  debtString.String(),
				Inline: false,
			},
		}
	}
	log.Debug().Msgf("transformed %v players to discord embed", len(players))

	return embed
}

func defaultEmbed() discord.Embed {
	version := os.Getenv("VERSION")
	return discord.Embed{
		Title:       ":moneybag: 10k in die Gildenbank!",
		Type:        discord.NormalEmbed,
		Description: "[GitHub](https://github.com/torfstack/slash10k) | v" + version,
		Timestamp:   discord.NowTimestamp(),
		Color:       discord.Color(0xF1C40F),
		Footer: &discord.EmbedFooter{
			Text: "/10k <Spieler> <Betrag>",
			Icon: "https://true.torfstack.com/coin.png",
		},
	}
}
