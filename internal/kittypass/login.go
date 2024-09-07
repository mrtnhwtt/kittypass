package kittypass

import (
	"encoding/hex"
	"log"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
)

type Login struct {
	Vault           *Vault
	Password        string
	Username        string
	Name            string
	HexSalt         string
	Salt            []byte
	DerivationKey   []byte
	ProvidePassword bool
	Generator       PasswordGenerator
}

func NewLogin() Login {
	return Login{
		Generator: PasswordGenerator{},
	}
}

// Create an hex of the salt and encrypted password 
func (l *Login) PrepareHex() {}

// UseMasterPassword generate a DerivationKey from a master password to use to encrypt and decrypt passwords stored in kittypass
func (l *Login) UseMasterPassword() error {
	var err error
	l.Salt, err = crypto.GenerateRandomSalt(16)
	if err != nil {
		return err
	}
	l.HexSalt = hex.EncodeToString(l.Salt)
	l.DerivationKey = crypto.GenerateKey([]byte(l.Vault.Masterpass), l.Salt)
	return nil
}

func (l *Login) RecreateDerivationKey() error {
	salt, err := hex.DecodeString(l.HexSalt)
	if err != nil {
		log.Printf("error when decoding salt: %s", err)
		return MalformedDataError{Data: "salt"}
	}
	l.DerivationKey = crypto.GenerateKey([]byte(l.Vault.Masterpass), salt)
	return nil
}

func (l *Login) Add() error {
	err := l.UseMasterPassword()
	if err != nil {
		return err
	}
	e := crypto.New("aes")
	cipher, err := e.Encrypt(l.DerivationKey, l.Password)
	if err != nil {
		return err
	}

	_, err = l.Vault.Db.SaveLogin(l.Vault.Uuid, l.Name, l.Username, cipher, l.HexSalt)
	if err != nil {
		return err
	}
	return nil
}

func (l *Login) Get() (map[string]string, error) {
	stored, err := l.Vault.Db.ReadLogin(l.Vault.Uuid, l.Name)
	if err != nil {
		return nil, err
	}
	l.HexSalt = stored["hex_salt"]
	err = l.RecreateDerivationKey()
	if err != nil {
		return nil, err
	}
	e := crypto.New("aes")
	decrypted, err := e.Decrypt(l.DerivationKey, stored["hex_encrypted_password"])
	if err != nil {
		return nil, err
	}
	login := map[string]string{
		"name":     stored["name"],
		"username": stored["username"],
		"password": decrypted,
	}
	return login, nil
}

func (l *Login) List() ([]map[string]string, error) {
	loginList, err := l.Vault.Db.ListLogin(l.Vault.Uuid, l.Name, l.Username)
	if err != nil {
		return nil, err
	}
	return loginList, nil
}

func (l *Login) Delete() error {
	return l.Vault.Db.DeleteLogin(l.Vault.Uuid, l.Name)
}

func (l *Login) Update(target string) (int64, error) {
	var cipher string
	var err error
	if l.Password != "" {
		err = l.UseMasterPassword()
		if err != nil {
			return 0, err
		}
		e := crypto.New("aes")
		cipher, err = e.Encrypt(l.DerivationKey, l.Password)
		if err != nil {
			return 0, err
		}
	}
	aff, err := l.Vault.Db.UpdateLogin(l.Vault.Uuid, target, l.Name, l.Username, cipher, l.HexSalt)
	if err != nil {
		return 0, err
	}
	return aff, nil
}
