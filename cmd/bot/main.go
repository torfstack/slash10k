package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"scurvy10k/internal/models"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
	{Name: "10kup", Description: "Setze den Channel für 10k updates", Options: discord.CommandOptions{
		&discord.ChannelOption{OptionName: "channel_id", Description: "Channel für 10k updates", Required: true},
	}},
}

var baseUrl = "https://true.torfstack.com/"
var debtsUrl = baseUrl + "api/debt"

var channelId discord.ChannelID
var messageId discord.MessageID

func main() {
	setupLogger()

	r := cmdroute.NewRouter()

	s := state.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentMessageContent)

	r.AddFunc("10ks", func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		debts, err := getDebts()
		if err != nil {
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + err.Error())}
		}
		return &api.InteractionResponseData{
			Embeds: &[]discord.Embed{
				*transformDebtsToEmbed(debts),
			},
		}
	})

	r.AddFunc("10k", func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()
		amount := options.Find("amount").String()

		res, err := http.Post(debtsUrl+"/"+name+"/"+amount, "application/json", nil)
		if err != nil {
			log.Error().Msgf("cannot post debt: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + err.Error())}
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Error().Msgf("cannot post debt: %s", res.Status)
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + res.Status)}
		}
		if channelId != discord.NullChannelID && messageId != discord.NullMessageID {
			updateDebtsMessage(s)
		}
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Added %v to %v", amount, name)),
			Flags:   discord.EphemeralMessage,
		}
	})

	r.AddFunc("10kup", func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		var err error
		cId, err := options.Find("channel_id").SnowflakeValue()
		if err != nil {
			log.Error().Msgf("cannot get channel_id: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Could not set channel")}
		}
		channelId = discord.ChannelID(cId)
		debts, err := getDebts()
		if err != nil {
			log.Error().Msgf("cannot get debts: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Could not get debts")}
		}
		m, err := s.SendMessage(discord.ChannelID(channelId), "", *transformDebtsToEmbed(debts))
		if err != nil {
			log.Error().Msgf("cannot send message: %s", err)
			return &api.InteractionResponseData{Content: option.NewNullableString("Could not send message")}
		}
		messageId = m.ID
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Channel set successfully"),
			Flags:   discord.EphemeralMessage,
		}
	})

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

func getDebts() (*models.AllDebtsResponse, error) {
	res, err := http.Get(debtsUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var debts models.AllDebtsResponse
	if err = json.Unmarshal(bytes, &debts); err != nil {
		log.Error().Msgf("cannot unmarshal debts: %s", err)
		return nil, err
	}
	return &debts, nil
}

func transformDebtsToEmbed(debts *models.AllDebtsResponse) *discord.Embed {
	var fields []discord.EmbedField
	for _, d := range debts.Debts {
		fields = append(fields, discord.EmbedField{
			Name:   d.Name,
			Value:  d.Amount,
			Inline: true,
		})
	}
	return &discord.Embed{
		Title:       "True",
		Type:        discord.NormalEmbed,
		Description: "10k in die Gildenbank!",
		URL:         baseUrl,
		Timestamp:   discord.NowTimestamp(),
		Color:       discord.DefaultEmbedColor,
		Footer: &discord.EmbedFooter{
			Text: "https://github.com/torfstack/scurvy10k",
			Icon: "https://github.githubassets.com/assets/GitHub-Mark-ea2971cee799.png",
		},
		Fields: fields,
	}
}

func updateDebtsMessage(s *state.State) {
	debts, err := getDebts()
	if err != nil {
		log.Error().Msgf("cannot get debts: %s", err)
		return
	}

	m, err := s.EditMessage(channelId, messageId, "", *transformDebtsToEmbed(debts))
	if err != nil {
		log.Error().Msgf("cannot edit message: %s", err)
		return
	}
	messageId = m.ID
}
