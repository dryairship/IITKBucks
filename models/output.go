package models

import (
	"encoding/binary"
)

type Output struct {
	Recipient User  `json:"recipient"`
	Amount    Coins `json:"amount"`
}

type OutputList []Output

type OutputMap = map[TransactionIdIndexPair]Output

func (output Output) ToByteArray() []byte {
	var result []byte

	amtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amtBytes, output.Amount)
	result = append(result, amtBytes...)

	rcptByteArray := output.Recipient.ToByteArray()
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(rcptByteArray)))
	result = append(result, lenBytes...)

	result = append(result, rcptByteArray...)

	return result
}

func (outputList OutputList) ToByteArray() []byte {
	var result []byte

	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(outputList)))
	result = append(result, lenBytes...)

	for i := range outputList {
		result = append(result, outputList[i].ToByteArray()...)
	}

	return result
}

func (outputList OutputList) GetSumOfAmounts() Coins {
	var totalCoins Coins = 0
	for i := range outputList {
		totalCoins += outputList[i].Amount
	}
	return totalCoins
}
