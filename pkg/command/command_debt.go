package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"os"
	"slash10k/pkg/domain"
	"slash10k/pkg/models"
	"slices"
	"strings"
)

func AddDebt(discord *state.State, service domain.Service) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()
		amount, err := options.Find("amount").IntValue()
		if err != nil || amount < 0 {
			log.Error().Msgf("could not parse amount: %v, err: %s", amount, err)
			return ephemeralMessage("Amount needs to be a non-negative number!")
		}
		reason := options.Find("reason").String()

		err = service.AddDebt(ctx, data.Event.User.ID.String(), data.Event.GuildID.String(), amount, reason)
		if err != nil {
			log.Error().Msgf("could not add debt: %s", err)
			return ephemeralMessage("Could not update debt")
		}

		updateDebtsMessage(ctx, discord, service, data.Event.GuildID.String())

		return ephemeralMessage(fmt.Sprintf("Added %v to %v, because '%v'", amount, name, reason))
	}
}

func SubDebt(discord *state.State, service domain.Service) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()
		amount, err := options.Find("amount").IntValue()
		if err != nil || amount < 0 {
			log.Error().Msgf("could not parse amount: %v, err: %s", amount, err)
			return ephemeralMessage("Amount needs to be a non-negative number!")
		}

		err = service.SubDebt(ctx, data.Event.User.ID.String(), data.Event.GuildID.String(), amount)
		if err != nil {
			log.Error().Msgf("could not subtract debt: %s", err)
			return ephemeralMessage("Could not update debt")
		}

		updateDebtsMessage(ctx, discord, service, data.Event.GuildID.String())

		return ephemeralMessage(fmt.Sprintf("Removed %v from %v", amount, name))
	}
}

func transformDebtsToEmbed(debts *models.AllDebtsResponse) *discord.Embed {
	embed := defaultEmbed()

	if len(debts.Debts) > 0 {
		maxLength := len(
			slices.MaxFunc(
				debts.Debts, func(d1, d2 models.PlayerDebt) int {
					return len(d1.Name) - len(d2.Name)
				},
			).Name,
		)
		debtString := strings.Builder{}
		debtString.WriteString("```")
		for _, d := range debts.Debts {
			debtString.WriteString(fmt.Sprintf("%-*s %v\n", maxLength, d.Name, d.Amount))
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
