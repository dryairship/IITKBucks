package models

type User string

func (user User) ToByteArray() []byte {
	return []byte(user)
}
