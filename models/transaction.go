package models

import (
	"crypto/sha256"
)

type Transaction struct {
	Id      ID         `json:"id"`
	Inputs  InputList  `json:"inputs"`
	Outputs OutputList `json:"outputs"`
}

type TransactionList []Transaction

func (txn *Transaction) AddInput(input Input) {
	txn.Inputs = append(txn.Inputs, input)
}

func (txn *Transaction) AddOutput(output Output) {
	txn.Outputs = append(txn.Outputs, output)
}

func (txn Transaction) ToByteArray() []byte {
	return append(txn.Inputs.ToByteArray(), txn.Outputs.ToByteArray()...)
}

func (txn Transaction) CalculateHash() Hash {
	hash := sha256.Sum256(txn.ToByteArray())
	return hash
}

func (txn Transaction) CalculateOutputDataHash() Hash {
	return sha256.Sum256(txn.Outputs.ToByteArray())
}
