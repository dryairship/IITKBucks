package models

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"time"
)

type Block struct {
	Id           ID              `json:"_"`
	Index        uint32          `json:"index"`
	Timestamp    int64           `json:"timestamp"`
	Transactions TransactionList `json:"transactions"`
	ParentHash   Hash            `json:"parentHash"`
	Target       Hash            `json:"target"`
	Nonce        int64           `json:"nonce"`
}

func NewGenesisBlock() *Block {
	return &Block{
		Index:     0,
		Timestamp: time.Now().UnixNano(),
	}
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

func BlockFromByteArray(data []byte) (Block, error) {
	if len(data) < 116 {
		return Block{}, ERROR_INSUFFICIENT_DATA
	}

	index := binary.BigEndian.Uint32(data[:4])

	var parentHash, blockHash, target Hash
	copy(parentHash[:], data[4:36])
	copy(blockHash[:], data[36:68])
	copy(target[:], data[68:100])

	timestamp := int64(binary.BigEndian.Uint64(data[100:108]))
	nonce := int64(binary.BigEndian.Uint64(data[108:116]))

	transactions, err := TransactionListFromByteArray(data[116:])
	if err != nil {
		return Block{}, err
	}

	return Block{
		Id:           blockHash,
		Index:        index,
		Timestamp:    timestamp,
		Transactions: transactions,
		ParentHash:   parentHash,
		Target:       target,
		Nonce:        nonce,
	}, nil
}
