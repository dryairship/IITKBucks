package models

type ID = Hash

type TransactionIdIndexPair struct {
	TransactionId ID
	Index         uint32
}
