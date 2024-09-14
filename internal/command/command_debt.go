package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"slash10k/internal/db"
	"slash10k/internal/models"
	"slices"
	"strings"
	"time"
)

func AddDebt(s *state.State, d db.Database, c DebtClient) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()
		amount, err := options.Find("amount").IntValue()
		if err != nil || amount < 0 {
			return ephemeralMessage("Amount needs to be a non-negative number!")
		}
		reason := options.Find("reason").String()

		err = c.AddDebt(ctx, name, amount, reason)
		if err != nil {
			log.Error().Msgf("could not add debt: %s", err)
			return ephemeralMessage("Could not update debt")
		}

		updateDebtsMessage(ctx, s, d, c)

		caser := cases.Title(language.English)
		return visibleMessage(fmt.Sprintf("Added %v to %v, because '%v'", amount, caser.String(name), reason))
	}
}

func SubDebt(s *state.State, d db.Database, c DebtClient) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()
		amount, err := options.Find("amount").IntValue()
		if err != nil || amount < 0 {
			return ephemeralMessage("Amount needs to be a non-negative number!")
		}

		err = c.AddDebt(ctx, name, -1*amount, "")
		if err != nil {
			log.Error().Msgf("could not subtract debt: %s", err)
			return ephemeralMessage("Could not update debt")
		}

		updateDebtsMessage(ctx, s, d, c)

		caser := cases.Title(language.English)
		return visibleMessage(fmt.Sprintf("Removed %v from %v", amount, caser.String(name)))
	}
}

func GetJournalEntries(c DebtClient) func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()

		entries, err := c.GetJournalEntries(ctx, name)

		if len(entries.Entries) > 0 {
			var s string
			s, err = journalEntriesString(*entries)
			if err != nil {
				log.Error().Msgf("cannot get journal entries string: %s", err)
				return ephemeralMessage("Could not get journal entries")
			}
			caser := cases.Title(language.English)
			return visibleMessage(fmt.Sprintf("Journal entries of %v%v", caser.String(name), s))
		} else {
			return ephemeralMessage("No journal entries found")
		}
	}
}

func journalEntriesString(entries models.JournalEntries) (string, error) {
	maxAmountLength := len(
		fmt.Sprint(
			slices.MaxFunc(
				entries.Entries, func(e1, e2 models.JournalEntry) int {
					return len(fmt.Sprint(e1.Amount)) - len(fmt.Sprint(e2.Amount))
				},
			).Amount,
		),
	)
	var b strings.Builder
	b.WriteString("```")
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return "", err
	}
	for _, entry := range entries.Entries {
		date := time.Unix(entry.Date, 0).In(berlin).Format(time.DateTime)
		b.WriteString(fmt.Sprintf("%-*v || %s || %v\n", maxAmountLength, entry.Amount, entry.Reason, date))
	}
	b.WriteString("```")
	return b.String(), nil
}

func RefreshDebts(s *state.State, d db.Database, c DebtClient) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		updateDebtsMessage(ctx, s, d, c)
		return ephemeralMessage("Debts refreshed successfully")
	}
}

func transformDebtsToEmbed(debts *models.AllDebtsResponse) *discord.Embed {
	maxLength := len(
		slices.MaxFunc(
			debts.Debts, func(d1, d2 models.PlayerDebt) int {
				return len(d1.Name) - len(d2.Name)
			},
		).Name,
	)
	debtString := strings.Builder{}
	debtString.WriteString("```")
	caser := cases.Title(language.English)
	for _, d := range debts.Debts {
		debtString.WriteString(fmt.Sprintf("%-*s %v\n", maxLength, caser.String(d.Name), d.Amount))
	}
	debtString.WriteString("```")

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
		Fields: []discord.EmbedField{
			{
				Name:   "Spieler",
				Value:  debtString.String(),
				Inline: false,
			},
		},
	}
}
