package models

type Output struct {
	Id        ID    `json:"id"`
	Recipient User  `json:"recipient"`
	Amount    Coins `json:"amount"`
}
