package main

import (
	"context"
	"encoding/json"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"scurvy10k/internal/models"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var commands = []api.CreateCommandData{
	{Name: "10ks", Description: "Wer packt 10k in die Gildenbank?"},
	{Name: "10k", Description: "Packt 10k in die Gildenbank!", Options: discord.CommandOptions{
		&discord.StringOption{OptionName: "name", Description: "Name des Spielers", Required: true},
		&discord.StringOption{OptionName: "amount", Description: "Betrag, kann negativ sein", Required: true},
	}},
}

var baseUrl string = "https://true.torfstack.com/"
var debtsUrl string = baseUrl + "api/debt"

func main() {
	setupLogger()

	r := cmdroute.NewRouter()

	r.AddFunc("10ks", func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		res, err := http.Get(debtsUrl)
		if err != nil {
			log.Error().Msgf("cannot get debts: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + err.Error())}
		}
		defer res.Body.Close()
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Msgf("cannot read body: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + err.Error())}
		}
		var debts models.AllDebtsResponse
		if err = json.Unmarshal(bytes, &debts); err != nil {
			log.Error().Msgf("cannot unmarshal debts: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + err.Error())}
		}
		b := strings.Builder{}
		for _, d := range debts.Debts {
			b.WriteString(d.Name + ": " + d.Amount + "\n")
		}
		return &api.InteractionResponseData{Content: option.NewNullableString(b.String())}
	})

	r.AddFunc("10k", func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		nameOption := options.Find("name")
		amount := options.Find("amount")

		res, err := http.Post(debtsUrl+"/"+nameOption.String()+"/"+amount.String(), "application/json", nil)
		if err != nil {
			log.Error().Msgf("cannot post debt: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + err.Error())}
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Error().Msgf("cannot post debt: %s", res.Status)
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + res.Status)}
		}
		return &api.InteractionResponseData{Content: option.NewNullableString("Erfolgreich!")}
	})

	s := state.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentMessageContent)

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
