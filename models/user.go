package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type User string

func (user User) ToByteArray() []byte {
	return []byte(user)
}

func (user User) GetPublicKey() (*rsa.PublicKey, error) {
	data, _ := pem.Decode(user.ToByteArray())
	if data == nil {
		return nil, ERROR_INVALID_PEM_BLOCK
	}

	pub, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}
	return nil, ERROR_NOT_RSA_KEY
}
