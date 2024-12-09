package domain

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
