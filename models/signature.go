package models

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
)

type Signature []byte

var unlockOptions = rsa.PSSOptions{
	SaltLength: 32,
	Hash:       crypto.SHA256,
}

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

func (sig Signature) Unlock(output *Output, message *Hash) bool {
	pubkey, err := output.Recipient.GetPublicKey()
	if err != nil {
		return false
	}

	if err := rsa.VerifyPSS(
		pubkey,
		crypto.SHA256,
		message.ToByteArray(),
		sig,
		&unlockOptions,
	); err != nil {
		return false
	}

	return true
}
