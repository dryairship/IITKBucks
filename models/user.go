package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/logger"
)

type User string

var MyUser = User(config.MY_PUBLIC_KEY)

func (user User) ToByteArray() []byte {
	return []byte(user)
}

func (user User) GetPublicKey() (*rsa.PublicKey, error) {
	data, _ := pem.Decode(user.ToByteArray())
	if data == nil {
		logger.Println(logger.RareError, "[Models/User] [ERROR] Invalid PEM block. User:", user)
		return nil, ERROR_INVALID_PEM_BLOCK
	}

	pub, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		logger.Println(logger.RareError, "[Models/User] [ERROR] Could not parse Public Key. User:", user, ", Err:", err)
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}

	logger.Println(logger.RareError, "[Models/User] [ERROR] User's key is not RSA Key. User:", user)
	return nil, ERROR_NOT_RSA_KEY
}
