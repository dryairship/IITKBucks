package models

import (
	"crypto/sha256"
	"encoding/json"
)

type Block struct {
	Index        int64         `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	ParentHash   Hash          `json:"parentHash"`
	Nonce        int64         `json:"nonce"`
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
