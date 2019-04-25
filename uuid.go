package maru

import (
	"crypto/rand"
	"fmt"
)

func GetUniqueId() string {
	var u = new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		panic(fmt.Sprintf("rand.Read error: %s",err))
	}
	return fmt.Sprintf("%x", u[:])
}