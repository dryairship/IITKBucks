package models

type Output struct {
	Id        ID
	Recipient User
	Amount    Coins
}
