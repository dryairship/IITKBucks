package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Hash [32]byte

func (hash Hash) String() string {
	return fmt.Sprintf("%x", [32]byte(hash))
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
