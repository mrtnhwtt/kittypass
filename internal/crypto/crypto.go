package crypto

import (
	"golang.org/x/crypto/argon2"
)

type Algorithm string

const (
	Ed25519 Algorithm = "ed25519"
)

type Encryption interface {
	encrypt()
	decrypt()
}

func GenerateKey() {
	return argon2.IDKey()
}
