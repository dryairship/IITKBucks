package models

type Transaction struct {
	Id      Hash       `json:"id"`
	Inputs  InputList  `json:"inputs"`
	Outputs OutputList `json:"outputs"`
}

func (txn Transaction) ToByteArray() []byte {
	return append(txn.Inputs.ToByteArray(), txn.Outputs.ToByteArray()...)
}
