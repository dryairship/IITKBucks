package models

import (
	"encoding/binary"
)

type ID = Hash

type TransactionIdIndexPair struct {
	TransactionId ID
	Index         uint32
}

func (pair TransactionIdIndexPair) ToByteArray() []byte {
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, uint32(pair.Index))
	return append(pair.TransactionId[:], indexBytes...)
}
