package models

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/dryairship/IITKBucks/logger"
)

type Signature []byte

var unlockOptions = rsa.PSSOptions{
	SaltLength: 32,
	Hash:       crypto.SHA256,
}

func (sig Signature) String() string {
	return fmt.Sprintf("%x", []byte(sig))
}

func (sig Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(sig.String())
}

func (sig Signature) ToByteArray() []byte {
	return sig
}

func SignatureFromHexString(str string) (Signature, error) {
	if str == "" {
		logger.Println(logger.RareError, "[Models/Signature] [ERROR] Signature is empty")
		return nil, ERROR_EMPTY_SIGNATURE_STRING
	}

	signature, err := hex.DecodeString(str)
	if err != nil {
		logger.Println(logger.RareError, "[Models/Signature] [ERROR] Could not decode signature from hex. Signature:", str)
		return nil, err
	}
	return Signature(signature), nil
}

func (sig Signature) Unlock(output *Output, txidIndexPair *TransactionIdIndexPair, message *Hash) bool {
	pubkey, err := output.Recipient.GetPublicKey()
	if err != nil {
		return false
	}

	totalData := append(txidIndexPair.ToByteArray(), message[:]...)
	hash := sha256.Sum256(totalData)

	if err := rsa.VerifyPSS(
		pubkey,
		crypto.SHA256,
		hash[:],
		sig,
		&unlockOptions,
	); err != nil {
		logger.Println(logger.RareError, "[Models/Signature] [ERROR] Transaction Verification Failed. TotalData:",
			hex.EncodeToString(totalData), ", Signature:", sig)
		return false
	}

	return true
}
