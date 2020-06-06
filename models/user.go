package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type User string

func (user User) ToByteArray() []byte {
	return []byte(user)
}

func (user User) GetPublicKey() (*rsa.PublicKey, error) {
	data, _ := pem.Decode(user.ToByteArray())
	if data == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
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
	return nil, errors.New("Key type is not RSA")
}
