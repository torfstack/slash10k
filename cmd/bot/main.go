package main

import (
	"context"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"slash10k/internal/command"
	"slash10k/internal/db"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var commands = []api.CreateCommandData{
	{
		Name: "10kup", Description: "Setze den Channel in dem der Bot aktiv sein soll", Options: discord.CommandOptions{
			&discord.ChannelOption{
				OptionName:  "channel_id",
				Description: "Channel, in dem der Bot aktiv sein soll",
				Required:    true,
			},
		},
	},
	{
		Name: "10k", Description: "Packt 10k in die Gildenbank!", Options: discord.CommandOptions{
			&discord.StringOption{OptionName: "name", Description: "Name des Spielers", Required: true},
			&discord.StringOption{OptionName: "amount", Description: "Betrag", Required: true},
			&discord.StringOption{OptionName: "reason", Description: "Grund", Required: true},
		},
	},
	{
		Name: "10kpay", Description: "Hat 10k in die Gildenbank gepackt!", Options: discord.CommandOptions{
			&discord.StringOption{OptionName: "name", Description: "Name des Spielers", Required: true},
			&discord.StringOption{OptionName: "amount", Description: "Betrag", Required: true},
		},
	},
	{
		Name: "10kwhy", Description: "Warum 10k? (Historie ist limitiert)", Options: discord.CommandOptions{
			&discord.StringOption{OptionName: "name", Description: "Name des Spielers", Required: true},
		},
	},
	{
		Name: "10kplayeradd", Description: "Füge einen Spieler hinzu.", Options: discord.CommandOptions{
			&discord.StringOption{OptionName: "name", Description: "Name des Spielers", Required: true},
		},
	},
	{
		Name: "10kplayerdel", Description: "Entferne einen Spieler.", Options: discord.CommandOptions{
			&discord.StringOption{OptionName: "name", Description: "Name des Spielers", Required: true},
		},
	},
	{Name: "10krefresh", Description: "Refresh debts."},
}

func main() {
	setupLogger()

	r := cmdroute.NewRouter()

	s := state.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentMessageContent)

	err := command.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup command")
	}

	d := db.NewDatabase()
	c := command.NewClient()

	r.AddFunc("10kup", command.SetChannel(s, d, c))
	r.AddFunc("10k", command.AddDebt(s, d, c))
	r.AddFunc("10kpay", command.SubDebt(s, d, c))
	r.AddFunc("10kwhy", command.GetJournalEntries(c))
	r.AddFunc("10kplayeradd", command.AddPlayer(s, d, c))
	r.AddFunc("10kplayerdel", command.DeletePlayer(s, d, c))
	r.AddFunc("10krefresh", command.RefreshDebts(s, d, c))

	if err := cmdroute.OverwriteCommands(s, commands); err != nil {
		log.Fatal().Msgf("cannot update commands: %s", err)
	}

	log.Info().Msg("connecting slash10k-bot")
	if err := s.Connect(context.Background()); err != nil {
		log.Fatal().Msgf("cannot connect: %s", err)
	}
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		l, err := zerolog.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else {
			zerolog.SetGlobalLevel(l)
		}
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
