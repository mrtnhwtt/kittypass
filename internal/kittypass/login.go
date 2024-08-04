package kittypass

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
	"github.com/mrtnhwtt/kittypass/internal/storage"
)

type Login struct {
	Vault           Vault
	Password        string
	Username        string
	Name            string
	ProvidePassword bool
	Generator       PasswordGenerator
}

type PasswordGenerator struct {
	Length      int
	SpecialChar bool
	Numeral     bool
	Uppercase   bool
}

func NewLogin() Login {
	return Login{
		Generator: PasswordGenerator{},
	}
}

func (l *Login) Add() error {
	e := crypto.New("aes")
	cipher, err := e.Encrypt(l.Vault.DerivationKey, l.Password)
	if err != nil {
		return fmt.Errorf("error while encrypting password: %s", err)
	}
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}

	_, err = db.SaveLogin(l.Vault.Uuid, l.Name, l.Username, cipher)
	if err != nil {
		return fmt.Errorf("error while saving to database: %s", err)
	}
	return nil
}

func (l *Login) Get(master string) (map[string]string, error) {
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	stored, err := db.ReadLogin(l.Vault.Uuid, l.Name)
	if err != nil {
		return nil, fmt.Errorf("error reading stored password: %s", err)
	}
	saltHex := stored["salt"]
	err = l.Vault.RecreateDerivationKey(saltHex, master)
	if err != nil {
		return nil, fmt.Errorf("error while creating derivation key from stored salt and masterpassword: %s", err)
	}
	e := crypto.New("aes")
	decrypted, err := e.Decrypt(l.Vault.DerivationKey, stored["encrypted_password"])
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
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	loginList, err := db.ListLogin(l.Vault.Uuid, l.Name, l.Username)
	if err != nil {
		return nil, err
	}
	return loginList, nil
}

func (l *Login) Delete() error {
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}
	return db.DeleteLogin(l.Vault.Uuid, l.Name)
}

func (l *Login) Update(target string) error {
	var cipher string
	var err error
	if l.Password != "" {
		e := crypto.New("aes")
		cipher, err = e.Encrypt(l.Vault.DerivationKey, l.Password)
		if err != nil {
			return fmt.Errorf("error while encrypting password: %s", err)
		}
	}
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}
	hexSalt := hex.EncodeToString(l.Vault.Salt)
	db.Update(l.Vault.Uuid, target, l.Name, l.Username, cipher, hexSalt)
	return nil
}

func (g *PasswordGenerator) GeneratePassword() string { // need to garantie one of each type if selected
	numbers := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
    lowercase := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
    uppercase := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
    specialChars := []rune{'!', '#', '$', '%', '&', '*', '+',  '-', '?', '@', '^', '_', '~'}

	characterPool := []rune{}
	characterPool = append(characterPool, lowercase...)
	if g.SpecialChar {
		characterPool = append(characterPool, specialChars...)
	}

	if g.Numeral {
		characterPool = append(characterPool, numbers...)
	}

	if g.Uppercase {
		characterPool = append(characterPool, uppercase...)
	}
	var password string
	for i := 0; i < g.Length; i++ {
		randind := rand.Intn(len(characterPool))
		password = password + string(characterPool[randind])
	}

	return password
}
