package models

type PlayerDebt struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type AllDebtsResponse struct {
	Debts []PlayerDebt `json:"debts"`
}

type JournalEntries struct {
	Entries []JournalEntry `json:"entries"`
}

type JournalEntry struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
	Date   int64  `json:"date"`
}
