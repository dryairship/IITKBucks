package models

import (
	"encoding/hex"
	"fmt"
)

type Signature []byte

func (sig Signature) String() string {
	return fmt.Sprintf("%x", []byte(sig))
}

func (sig Signature) ToByteArray() []byte {
	return sig
}

func SignatureFromHexString(str string) (Signature, error) {
	signature, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return Signature(signature), nil
}
