package models

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
)

type Block struct {
	Id           ID              `json:"_"`
	Index        uint32          `json:"index"`
	ParentHash   Hash            `json:"parentHash"`
	BodyHash     Hash            `json:"bodyHash"`
	Target       Hash            `json:"target"`
	Timestamp    int64           `json:"timestamp"`
	Nonce        int64           `json:"nonce"`
	Transactions TransactionList `json:"transactions"`
}

func NewGenesisBlock() Block {
	output := Output{
		Amount:    Blockchain().CurrentBlockReward,
		Recipient: MyUser,
	}
	outputList := OutputList{output}
	transaction := Transaction{
		Outputs: outputList,
	}
	transaction.CalculateHash()
	transactionList := TransactionList{transaction}
	block := Block{
		Index:        0,
		Target:       Blockchain().CurrentTarget,
		Transactions: transactionList,
	}
	return block
}

func (block Block) ToJSON() string {
	json, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (block Block) GetHash() Hash {
	return sha256.Sum256(block.CalculateHeader(false))
}

func BlockFromByteArray(data []byte) (Block, error) {
	if len(data) < 116 {
		return Block{}, ERROR_INSUFFICIENT_DATA
	}

	index := binary.BigEndian.Uint32(data[:4])

	var parentHash, bodyHash, target Hash
	copy(parentHash[:], data[4:36])
	copy(bodyHash[:], data[36:68])
	copy(target[:], data[68:100])

	timestamp := int64(binary.BigEndian.Uint64(data[100:108]))
	nonce := int64(binary.BigEndian.Uint64(data[108:116]))

	transactions, err := TransactionListFromByteArray(data[116:])
	if err != nil {
		return Block{}, err
	}

	return Block{
		Index:        index,
		Timestamp:    timestamp,
		Transactions: transactions,
		ParentHash:   parentHash,
		BodyHash:     bodyHash,
		Target:       target,
		Nonce:        nonce,
	}, nil
}

func (block Block) GetBodyHash() Hash {
	return sha256.Sum256(block.Transactions.ToByteArray())
}

func (block *Block) CalculateHeader(recalculateBodyHash bool) []byte {
	var result []byte

	fourBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(fourBytes, block.Index)
	result = append(result, fourBytes...)

	result = append(result, block.ParentHash[:]...)

	if recalculateBodyHash {
		bodyHash := block.GetBodyHash()
		block.BodyHash = bodyHash
		result = append(result, bodyHash[:]...)
	} else {
		result = append(result, block.BodyHash[:]...)
	}

	result = append(result, block.Target[:]...)

	eightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(eightBytes, uint64(block.Timestamp))
	result = append(result, eightBytes...)

	binary.BigEndian.PutUint64(eightBytes, uint64(block.Nonce))
	result = append(result, eightBytes...)

	return result
}
