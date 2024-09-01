package crypto

import "fmt"

type EncryptionKeyError struct {
	Message string
}

func (e EncryptionKeyError) Error() string {
	return fmt.Sprintf("invalid encryption key: %s", e.Message)
}

type GenEncryptionKeyError struct {
	Message string
}

func (e GenEncryptionKeyError) Error() string {
	return fmt.Sprintf("failed to generate encryption key: %s", e.Message)
}

type EncryptionError struct {}

func (e EncryptionError) Error() string {
	return "failed to encrypt"
}

type DecryptionError struct {}

func (e DecryptionError) Error() string {
	return "failed to decrypt"
}

type MalformedDataError struct {
	Data string
}

func (e MalformedDataError) Error() string {
	return fmt.Sprintf("data for %s is malformed, could not be processed", e.Data)
}