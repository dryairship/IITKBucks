package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Hash [32]byte

func (hash Hash) String() string {
	return fmt.Sprintf("%x", [32]byte(hash))
}

func HashFromHexString(str string) (Hash, error) {
	var hash [32]byte

	hashSlice, err := hex.DecodeString(str)
	if err != nil {
		return [32]byte{}, err
	}

	copy(hash[:], hashSlice[:32])
	return Hash(hash), nil
}

func (hash Hash) ToByteArray() []byte {
	return hash[:]
}

func (hash Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(hash.String())
}

func (hash Hash) IsLessThan(target Hash) bool {
	return bytes.Compare(hash[:], target[:]) == -1
}
