package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scurvy10k/internal/models"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var baseUrl = "https://true.torfstack.com/"
var debtsUrl = baseUrl + "api/debt"

var channelId discord.ChannelID
var messageId discord.MessageID

func GetDebts(s *state.State) func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		debts, err := getDebts()
		if err != nil {
			return &api.InteractionResponseData{Content: option.NewNullableString("Error: " + err.Error())}
		}
		return &api.InteractionResponseData{
			Embeds: &[]discord.Embed{
				*transformDebtsToEmbed(debts),
			},
		}
	}
}

func AddDebt(s *state.State) func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
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
	}
}

func SetChannel(s *state.State) func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
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
	}
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
	maxLength := len(slices.MaxFunc(debts.Debts, func(d1, d2 models.PlayerDebt) int {
		return len(d1.Name) - len(d2.Name)
	}).Name)
	debtString := strings.Builder{}
	debtString.WriteString("```")
	for _, d := range debts.Debts {
		debtString.WriteString(fmt.Sprintf("%-*s %v\n", maxLength, d.Name, d.Amount))
	}
	debtString.WriteString("```")

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
		Fields: []discord.EmbedField{
			{
				Name:   "Spieler",
				Value:  debtString.String(),
				Inline: false,
			},
		},
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
