package models

import (
	"sort"
	"strings"
)

type Players []Player

func (p Players) SortByName() {
	sort.Slice(
		p, func(i, j int) bool {
			return strings.Compare(p[i].Name, p[j].Name) < 0
		},
	)
}

type Player struct {
	Id          int32
	DiscordId   string
	DiscordName string
	GuildId     string
	Name        string
	Debt        Debt
	DebtJournal []DebtJournalEntry
}

type Debt struct {
	Id          int32
	Amount      int64
	LastUpdated int64
	UserId      int32
	GuildId     string
}

type DebtJournalEntry struct {
	Id          int32
	Amount      int64
	Description string
	Date        int64
	UserId      int32
	GuildId     string
}

type BotSetup struct {
	GuildId               string
	ChannelId             string
	RegistrationMessageId string
	DebtsMessageId        string
}
