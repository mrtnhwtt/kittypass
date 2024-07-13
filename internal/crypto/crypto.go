package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/argon2"
)

type Encryption interface {
	Encrypt(masterKey []byte, plainText string) (string, error)
	Decrypt(masterKey []byte, cipherText string) (string, error)
}

func GenerateKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
}

func GenerateRandomSalt(saltSize int) []byte {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return salt
}

func New(t string) Encryption {
	switch t {
	case "aes":
		return Aes{}
	default:
		return Aes{}
	}

}

type Aes struct {
}

func (a Aes) Encrypt(masterKey []byte, plainText string) (string, error) {
	cb, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return hex.EncodeToString(cipherText), nil
}

func (a Aes) Decrypt(masterKey []byte, cipherText string) (string, error) {
	ct, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	cb, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		return "", err
	}

	out, err := gcm.Open(nil, ct[:gcm.NonceSize()], ct[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
