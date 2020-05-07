package models

type Transaction struct {
	Id      ID
	Input   Input
	Outputs []Output
}
