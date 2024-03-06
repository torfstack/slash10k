package models

type PlayerDebt struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type AllDebtsResponse struct {
	Debts []PlayerDebt `json:"debts"`
}
