package models

type Transaction struct {
	Id      ID       `json:"id"`
	Input   Input    `json:"input"`
	Outputs []Output `json:"outputs"`
}
