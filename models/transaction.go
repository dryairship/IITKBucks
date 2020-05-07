package models

type Transaction struct {
	Id      ID         `json:"id"`
	Input   ID         `json:"input"`
	Outputs OutputList `json:"outputs"`
}
