package models

import (
	"encoding/binary"
)

type Output struct {
	Recipient User  `json:"recipient"`
	Amount    Coins `json:"amount"`
}
type OutputList []Output

func (output Output) ToByteArray() []byte {
	var result []byte

	result = append(result, output.Recipient.ToByteArray()...)

	amtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amtBytes, output.Amount)
	result = append(result, amtBytes...)

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

func (outputList *OutputList) DeleteOutputAtIndex(i int) {
	(*outputList)[i] = (*outputList)[len(*outputList)-1]
	*outputList = (*outputList)[:len(*outputList)-1]
}

func (outputList OutputList) GetSumOfAmounts() Coins {
	var totalCoins Coins = 0
	for i := range outputList {
		totalCoins += outputList[i].Amount
	}
	return totalCoins
}
