package models

import (
	"fmt"
)

type Hash [32]byte

func (hash Hash) String() string {
	return fmt.Sprintf("%x", [32]byte(hash))
}
