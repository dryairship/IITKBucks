package models

import (
	"encoding/binary"
)

type Input struct {
	TransactionId Hash
	OutputIndex   uint32
	Signature     Signature
}

type InputList []Input

func (input Input) ToByteArray() []byte {
	var result []byte

	result = append(result, input.TransactionId.ToByteArray()...)

	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, input.OutputIndex)
	result = append(result, indexBytes...)

	result = append(result, input.Signature.ToByteArray()...)

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
