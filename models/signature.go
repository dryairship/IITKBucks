package models

import (
	"encoding/hex"
	"errors"
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
	if str == "" {
		return nil, errors.New("Empty string provided as signature")
	}

	signature, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return Signature(signature), nil
}
