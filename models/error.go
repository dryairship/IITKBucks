package models

import (
	"errors"
)

var ERROR_INSUFFICIENT_DATA = errors.New("Insufficient data")
var ERROR_INCORRECT_HASH_STRING_LENGTH = errors.New("string does not have the correct length to represent a hash")
var ERROR_EMPTY_SIGNATURE_STRING = errors.New("string is empty, cannot represent a signature")
var ERROR_NOT_RSA_KEY = errors.New("the key is not an RSA key")
var ERROR_INVALID_PEM_BLOCK = errors.New("PEM block is invalid")
var ERROR_INVALID_BLOCK = errors.New("Invalid Block")

var ERROR_INCORRECT_INDEX = errors.New("Block has incorrect index")
var ERROR_LESS_TIMESTAMP = errors.New("Block's timestamp is less than parent's timestamp")
var ERROR_PARENT_HASH_MISMATCH = errors.New("Parent hash doesn't match with the hash of the block at that index")
var ERROR_TARGET_MISMATCH = errors.New("Block's target doesn't match with the current target")
var ERROR_INCORRECT_BODY_HASH = errors.New("Block's body hash has been incorrectly calculated")
var ERROR_ABOVE_TARGET = errors.New("Block's hash is not below the target")
var ERROR_CONTAINS_INVALID_TRANSACTION = errors.New("Block contains invalid transactions")
var ERROR_INVALID_COINBASE_TRANSCTION = errors.New("Block's coinbase transaction collects more coins than expected")
