package kittypass

import (
	"fmt"
	"math/rand"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
	"github.com/mrtnhwtt/kittypass/internal/storage"
)

type Login struct {
	Vault           *Vault
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
	defer db.Close()

	_, err = db.SaveLogin(l.Vault.Uuid, l.Name, l.Username, cipher)
	if err != nil {
		return fmt.Errorf("error while saving to database: %s", err)
	}
	return nil
}

func (l *Login) Get() (map[string]string, error) {
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	defer db.Close()

	stored, err := db.ReadLogin(l.Vault.Uuid, l.Name)
	if err != nil {
		return nil, fmt.Errorf("error reading stored password: %s", err)
	}
	err = l.Vault.RecreateDerivationKey()
	if err != nil {
		return nil, fmt.Errorf("error while creating derivation key from stored salt and masterpassword: %s", err)
	}
	e := crypto.New("aes")
	decrypted, err := e.Decrypt(l.Vault.DerivationKey, stored["hex_encrypted_password"])
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
	defer db.Close()
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
	defer db.Close()
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
	defer db.Close()
	err = db.UpdateLogin(l.Vault.Uuid, target, l.Name, l.Username, cipher)
	if err != nil {
		return err
	}
	return nil
}

func (g *PasswordGenerator) GeneratePassword() string {
	lowercase := []rune("abcdefghijklmnopqrstuvwxyz")
	numbers := []rune("0123456789")
	uppercase := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialChars := []rune("!#$%&*+-?@^_~")

	var selectedChars []rune

	selectedChars = append(selectedChars, lowercase[rand.Intn(len(lowercase))])

	if g.Numeral {
		selectedChars = append(selectedChars, numbers[rand.Intn(len(numbers))])
	}
	if g.Uppercase {
		selectedChars = append(selectedChars, uppercase[rand.Intn(len(uppercase))])
	}
	if g.SpecialChar {
		selectedChars = append(selectedChars, specialChars[rand.Intn(len(specialChars))])
	}

	allChars := lowercase
	if g.Numeral {
		allChars = append(allChars, numbers...)
	}
	if g.Uppercase {
		allChars = append(allChars, uppercase...)
	}
	if g.SpecialChar {
		allChars = append(allChars, specialChars...)
	}

	remainingLength := g.Length - len(selectedChars)
	for i := 0; i < remainingLength; i++ {
		selectedChars = append(selectedChars, allChars[rand.Intn(len(allChars))])
	}

	rand.Shuffle(len(selectedChars), func(i, j int) {
		selectedChars[i], selectedChars[j] = selectedChars[j], selectedChars[i]
	})

	return string(selectedChars)
}
