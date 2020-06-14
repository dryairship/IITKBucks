package models

import (
	"encoding/binary"
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

func (inputRequestBody InputRequestBody) ToInput() (Input, error) {
	txid, err := HashFromHexString(inputRequestBody.TransactionId)
	if err != nil {
		return Input{}, err
	}

	sign, err := SignatureFromHexString(inputRequestBody.Signature)
	if err != nil {
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
	for _, inputRequestBody := range inputListRequestBody {
		input, err := inputRequestBody.ToInput()
		if err != nil {
			return inputList, err
		}
		inputList = append(inputList, input)
	}
	return inputList, nil
}
