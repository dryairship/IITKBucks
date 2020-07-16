package controllers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"

	"github.com/gin-gonic/gin"
)

var signOptions = &rsa.PSSOptions{
	SaltLength: 32,
	Hash:       crypto.SHA256,
}

type signRequest struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

func readKey(keyStr string) (*rsa.PrivateKey, error) {
	pemKey, _ := pem.Decode([]byte(keyStr))
	return x509.ParsePKCS1PrivateKey(pemKey.Bytes)
}

func handleSign(c *gin.Context) {
	var body signRequest
	err := c.BindJSON(&body)
	if err != nil {
		c.String(400, "Invalid JSON request body")
		return
	}

	key, err := readKey(body.Key)
	if err != nil {
		c.String(400, "Invalid key")
		return
	}

	data, err := hex.DecodeString(body.Data)
	if err != nil {
		c.String(400, "Invalid data")
		return
	}

	hash := sha256.Sum256(data)

	sign, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, hash[:], signOptions)
	if err != nil {
		c.String(400, "Could not create signature")
		return
	}

	c.JSON(200, gin.H{
		"signature": hex.EncodeToString(sign),
	})
}
