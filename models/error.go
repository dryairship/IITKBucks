package models

import (
	"errors"
)

var ERROR_INSUFFICIENT_DATA = errors.New("Insufficient data")
var ERROR_INCORRECT_HASH_STRING_LENGTH = errors.New("string does not have the correct length to represent a hash")
var ERROR_EMPTY_SIGNATURE_STRING = errors.New("string is empty, cannot represent a signature")
var ERROR_NOT_RSA_KEY = errors.New("the key is not an RSA key")
var ERROR_INVALID_PEM_BLOCK = errors.New("PEM block is invalid")
