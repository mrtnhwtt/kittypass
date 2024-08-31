package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"

	"golang.org/x/crypto/argon2"
)

type Encryption interface {
	Encrypt(masterKey []byte, plainText string) (string, error)
	Decrypt(masterKey []byte, cipherText string) (string, error)
}

func GenerateKey(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
}

func GenerateRandomSalt(saltSize int) ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("error while generating encryption key unique salt: %s", err)
		return []byte{}, GenEncryptionKeyError{Message: "error while generating encryption key unique salt"}
	}
	return salt, nil
}

func New(t string) Encryption {
	switch t {
	case "aes":
		return Aes{}
	default:
		return Aes{}
	}

}

type Aes struct{}

func (a Aes) Encrypt(masterKey []byte, plainText string) (string, error) {
	cb, err := aes.NewCipher(masterKey)
	if err != nil {
		if _, ok := err.(aes.KeySizeError); ok {
			log.Printf("encryption function received invalid encryption key length of %d, expect 32", len(masterKey))
			return "", EncryptionKeyError{Message: "invalid encryption key length"}
		}
		log.Printf("failed to generate cipher block for encryption: %s", err)
		return "", EncryptionError{}
	}

	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		log.Printf("failed to generate encryption algorithm from cipher block: %s", err)
		return "", EncryptionError{}
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("failed to fill nonce: %s", err)
		return "", EncryptionError{}
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return hex.EncodeToString(cipherText), nil
}

func (a Aes) Decrypt(masterKey []byte, cipherText string) (string, error) {
	ct, err := hex.DecodeString(cipherText)
	if err != nil {
		log.Printf("error when decoding stored hex password: %s", err)
		return "", MalformedDataError{Data: "hex encrypted password"}
	}

	cb, err := aes.NewCipher(masterKey)
	if err != nil {
		if _, ok := err.(aes.KeySizeError); ok {
			log.Printf("encryption function received invalid encryption key length of %d, expect 32", len(masterKey))
			return "", EncryptionKeyError{Message: "invalid encryption key length"}
		}
		log.Printf("failed to generate cipher block for encryption: %s", err)
		return "", EncryptionError{}
	}
	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		log.Printf("failed to generate nonce generator: %s", err)
		return "", DecryptionError{}
	}

	out, err := gcm.Open(nil, ct[:gcm.NonceSize()], ct[gcm.NonceSize():], nil)
	if err != nil {
		log.Printf("failed to decrypt cipher text: %s", err)
		return "", DecryptionError{}
	}
	return string(out), nil
}
