package command

import (
	"encoding/base64"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func ephemeralMessage(content string) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content: option.NewNullableString(content),
		Flags:   discord.EphemeralMessage,
	}
}

func visibleMessage(content string) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content: option.NewNullableString(content),
	}
}

func basicAuth(user, pass string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
}
