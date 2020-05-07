package models

import (
	"crypto/sha256"
	"encoding/json"
)

type Block struct {
	Index        int64
	Timestamp    Timestamp
	Transactions []Transaction
	ParentHash   Hash
	Nonce        int64
}

func (block Block) ToJSON() string {
	json, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (block Block) GetHash() Hash {
	bytes := []byte(block.ToJSON())
	sum := sha256.Sum256(bytes)
	return sum
}
