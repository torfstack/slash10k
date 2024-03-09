package main

import (
	"context"
	"os"
	"scurvy10k/internal/command"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var commands = []api.CreateCommandData{
	{Name: "10ks", Description: "Wer packt 10k in die Gildenbank?"},
	{Name: "10k", Description: "Packt 10k in die Gildenbank!", Options: discord.CommandOptions{
		&discord.StringOption{OptionName: "name", Description: "Name des Spielers", Required: true},
		&discord.StringOption{OptionName: "amount", Description: "Betrag, kann negativ sein", Required: true},
	}},
	{Name: "10kchannel", Description: "Setze den Channel für 10k updates", Options: discord.CommandOptions{
		&discord.ChannelOption{OptionName: "channel_id", Description: "Channel für 10k updates", Required: true},
	}},
}

func main() {
	setupLogger()

	r := cmdroute.NewRouter()

	s := state.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentMessageContent)

	r.AddFunc("10ks", command.GetDebts(s))
	r.AddFunc("10k", command.AddDebt(s))
	r.AddFunc("10kchannel", command.SetChannel(s))

	if err := cmdroute.OverwriteCommands(s, commands); err != nil {
		log.Fatal().Msgf("cannot update commands: %s", err)
	}

	log.Info().Msg("connecting scurvy10k-bot")
	if err := s.Connect(context.TODO()); err != nil {
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
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(output)
}
