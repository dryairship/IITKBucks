package models

import (
	"encoding/binary"

	"github.com/dryairship/IITKBucks/logger"
)

type Input struct {
	TransactionId Hash      `json:"transactionId"`
	OutputIndex   uint32    `json:"index"`
	Signature     Signature `json:"signature"`
}

type InputRequestBody struct {
	TransactionId string `json:"transactionId" binding:"required"`
	OutputIndex   uint32 `json:"index" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
}

type InputList []Input
type InputListRequestBody []InputRequestBody

func (input Input) ToByteArray() []byte {
	var result []byte

	result = append(result, input.TransactionId.ToByteArray()...)

	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, input.OutputIndex)
	result = append(result, indexBytes...)

	sgnByteArray := input.Signature.ToByteArray()
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(sgnByteArray)))
	result = append(result, lenBytes...)

	result = append(result, sgnByteArray...)

	return result
}

func (inputList InputList) ToByteArray() []byte {
	var result []byte

	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(inputList)))
	result = append(result, lenBytes...)

	for i := range inputList {
		result = append(result, inputList[i].ToByteArray()...)
	}

	return result
}

func InputFromByteArray(data []byte) (Input, int, error) {
	if len(data) < 40 {
		logger.Println(logger.RareError, "[Models/Input] [ERROR] Input has insufficient data")
		return Input{}, 0, ERROR_INSUFFICIENT_DATA
	}

	var txnID Hash
	copy(txnID[:], data[:32])

	outputIndex := binary.BigEndian.Uint32(data[32:36])
	sigSize := int(binary.BigEndian.Uint32(data[36:40]))

	if len(data) < sigSize+40 {
		logger.Println(logger.RareError, "[Models/Input] [ERROR] Input has insufficient data")
		return Input{}, 0, ERROR_INSUFFICIENT_DATA
	}

	sign := Signature(data[40:sigSize])

	return Input{
		TransactionId: txnID,
		OutputIndex:   outputIndex,
		Signature:     sign,
	}, sigSize + 40, nil
}

func InputListFromByteArray(data []byte) (InputList, int, error) {
	if len(data) < 4 {
		logger.Println(logger.RareError, "[Models/Input] [ERROR] InputList has insufficient data")
		return nil, 0, ERROR_INSUFFICIENT_DATA
	}

	numInputs := binary.BigEndian.Uint32(data[:4])
	currentOffset := 4

	var list InputList
	for i := uint32(0); i < numInputs; i++ {
		input, bytesRead, err := InputFromByteArray(data[currentOffset:])
		if err != nil {
			logger.Printf(logger.RareError, "[Models/Input] [ERROR] Input %d has error. Err: %v\n", i, err)
			return list, 0, err
		}
		list = append(list, input)
		currentOffset += bytesRead
	}

	return list, currentOffset, nil
}

func (inputRequestBody InputRequestBody) ToInput() (Input, error) {
	txid, err := HashFromHexString(inputRequestBody.TransactionId)
	if err != nil {
		logger.Println(logger.RareError, "[Models/Input] [ERROR] Error while parsing Transaction ID from hex in input. Err:", err)
		return Input{}, err
	}

	sign, err := SignatureFromHexString(inputRequestBody.Signature)
	if err != nil {
		logger.Println(logger.RareError, "[Models/Input] [ERROR] Error while parsing Signature from hex in input. Err:", err)
		return Input{}, err
	}

	return Input{
		TransactionId: txid,
		Signature:     sign,
		OutputIndex:   inputRequestBody.OutputIndex,
	}, nil
}

func (inputListRequestBody InputListRequestBody) ToInputList() (InputList, error) {
	var inputList InputList
	for i, inputRequestBody := range inputListRequestBody {
		input, err := inputRequestBody.ToInput()
		if err != nil {
			logger.Printf(logger.RareError, "[Models/Input] [ERROR] Input %d has error. Err: %v\n", i, err)
			return inputList, err
		}
		inputList = append(inputList, input)
	}
	return inputList, nil
}
