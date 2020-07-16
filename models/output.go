package models

import (
	"encoding/binary"
	"encoding/json"

	"github.com/dryairship/IITKBucks/logger"
)

type Output struct {
	Recipient User  `json:"recipient"`
	Amount    Coins `json:"amount"`
}

type OutputRequestBody struct {
	Recipient string      `json:"recipient" binding:"required"`
	Amount    json.Number `json:"amount" binding:"required"`
}

type OutputList []Output
type OutputListRequestBody []OutputRequestBody

type OutputMap = map[TransactionIdIndexPair]Output

func (output Output) ToByteArray() []byte {
	var result []byte

	amtBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amtBytes, uint64(output.Amount))
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

func OutputFromByteArray(data []byte) (Output, int, error) {
	if len(data) < 12 {
		logger.Println(logger.RareError, "[Models/Output] [ERROR] Output has less than 12 bytes")
		return Output{}, 0, ERROR_INSUFFICIENT_DATA
	}

	amount := binary.BigEndian.Uint64(data[:8])
	pubKeySize := int(binary.BigEndian.Uint32(data[8:12]))

	if len(data) < pubKeySize+12 {
		logger.Println(logger.RareError, "[Models/Input] [ERROR] Output has insufficient data")
		return Output{}, 0, ERROR_INSUFFICIENT_DATA
	}

	recipient := User(data[12 : 12+pubKeySize])

	return Output{
		Amount:    Coins(amount),
		Recipient: recipient,
	}, pubKeySize + 12, nil
}

func OutputListFromByteArray(data []byte) (OutputList, int, error) {
	if len(data) < 4 {
		logger.Println(logger.RareError, "[Models/Output] [ERROR] OutputList has less than 4 bytes")
		return nil, 0, ERROR_INSUFFICIENT_DATA
	}

	numOutputs := binary.BigEndian.Uint32(data[:4])
	currentOffset := 4

	var list OutputList
	for i := uint32(0); i < numOutputs; i++ {
		output, bytesRead, err := OutputFromByteArray(data[currentOffset:])
		if err != nil {
			logger.Printf(logger.RareError, "[Models/Output] [ERROR] Output %d has error: %v\n", i, err)
			return list, 0, err
		}
		list = append(list, output)
		currentOffset += bytesRead
	}

	return list, currentOffset, nil
}

func (outputList OutputList) GetSumOfAmounts() Coins {
	var totalCoins Coins = 0
	for i := range outputList {
		totalCoins += outputList[i].Amount
	}
	return totalCoins
}

func (outputRequestBody OutputRequestBody) ToOutput() (Output, error) {
	amt, err := outputRequestBody.Amount.Int64()
	if err != nil {
		logger.Println(logger.RareError, "[Models/Output] [ERROR] Could not parse amounts from string to Coins. Err:", err)
		return Output{}, err
	}

	return Output{
		Recipient: User(outputRequestBody.Recipient),
		Amount:    Coins(amt),
	}, nil
}

func (outputListRequestBody OutputListRequestBody) ToOutputList() (OutputList, error) {
	var outputList OutputList
	for i, outputRequestBody := range outputListRequestBody {
		thisOutput, err := outputRequestBody.ToOutput()
		if err != nil {
			logger.Printf(logger.RareError, "[Models/Output] [ERROR] Output %d has error. Err: %v\n", i, err)
			return outputList, err
		}
		outputList = append(outputList, thisOutput)
	}
	return outputList, nil
}
