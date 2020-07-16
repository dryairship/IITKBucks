package models

import (
	"fmt"
)

type Coins uint64

func (coins Coins) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", uint64(coins))), nil
}
