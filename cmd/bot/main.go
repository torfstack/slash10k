package main

import (
	"context"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"slash10k/pkg/command"
	"slash10k/pkg/config"
	"slash10k/pkg/convert/fromdb"
	"slash10k/pkg/db"
	"slash10k/pkg/domain"
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
}

func main() {
	setupLogger()

	r := cmdroute.NewRouter()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal().Msg("DISCORD_TOKEN not set")
	}

	s := state.New("Bot " + token)
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentMessageContent)
	s.AddIntents(gateway.IntentGuildMessageReactions)

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get config from env")
	}
	err = db.Migrate(context.Background(), cfg.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("could not migrate database")
	}
	d := db.NewDatabase(cfg.ConnectionString())

	conn, err := d.Connect(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}
	botSetups, err := conn.Queries().GetAllBotSetups(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("could not get bot setups")
	}
	messageLookup := domain.NewMessageLookup(fromdb.FromBotSetups(botSetups))

	service := domain.NewSlashTenK(d)

	command.RegisterDiscordHandlers(s, service, messageLookup)

	r.AddFunc("10kup", command.SetChannel(s, service, messageLookup))

	if err := cmdroute.OverwriteCommands(s, commands); err != nil {
		log.Fatal().Msgf("cannot update commands: %s", err)
	}

	log.Info().Msg("connecting slash10k-bot")
	if err := s.Connect(context.Background()); err != nil {
		log.Fatal().Msgf("cannot connect: %s", err)
	}
}

func setupLogger() {
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
