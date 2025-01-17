package pkg

import (
	"crypto/sha1"
	"fmt"
)

const salt = "703f2h08g723fh"

func GenerateHashedPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
