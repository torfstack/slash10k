package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"os"
	"slash10k/pkg/domain"
	"slash10k/pkg/models"
	"slices"
	"strings"
	_ "time/tzdata"
)

func updateDebtsMessage(ctx context.Context, state *state.State, service domain.Service, guildId string) {
	allPlayers, err := service.GetAllPlayers(ctx, guildId)
	if err != nil {
		log.Error().Msgf("cannot get all players: %s", err)
		return
	}

	botSetup, err := service.GetBotSetup(ctx, guildId)
	if err != nil {
		log.Error().Msgf("cannot get bot setup: %s", err)
		return
	}
	channelId, messageId := botSetupToDiscordTypes(*botSetup)

	if channelId == discord.NullChannelID || messageId == discord.NullMessageID {
		log.Error().Msg("channel id or message id is null")
		return
	}

	_, err = state.EditMessage(
		channelId,
		messageId,
		"",
		*transformDebtsToEmbed(allPlayers),
	)
	if err != nil {
		log.Error().Msgf("cannot edit message: %s", err)
	}
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

func transformDebtsToEmbed(players []models.Player) *discord.Embed {
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

	return embed
}

func defaultEmbed() *discord.Embed {
	version := os.Getenv("VERSION")
	return &discord.Embed{
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
