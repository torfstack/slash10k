package config

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/rs/zerolog/log"
)

type DiscordGuildId string

type DiscordUserId string

const (
	TorfstackServerGuildIdString DiscordGuildId = "309323862326116352"
	TrueServerGuildIdString      DiscordGuildId = "822223333030101092"

	TorfstackUserIdString DiscordUserId = "263352209654153236"
)

func TorfstackUserId() discord.UserID {
	id, err := discord.ParseSnowflake(string(TorfstackUserIdString))
	if err != nil {
		log.Err(err).Msg("could not convert torfstack-id to snowflake")
		return 0
	}
	return discord.UserID(id)
}
