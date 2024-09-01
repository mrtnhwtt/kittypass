package kittypass

import (

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

func NewLogin() Login {
	return Login{
		Generator: PasswordGenerator{},
	}
}

func (l *Login) Add() error {
	e := crypto.New("aes")
	cipher, err := e.Encrypt(l.Vault.DerivationKey, l.Password)
	if err != nil {
		return err
	}
	db, err := storage.New("./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.SaveLogin(l.Vault.Uuid, l.Name, l.Username, cipher)
	if err != nil {
		return err
	}
	return nil
}

func (l *Login) Get() (map[string]string, error) {
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stored, err := db.ReadLogin(l.Vault.Uuid, l.Name)
	if err != nil {
		return nil, err
	}
	err = l.Vault.RecreateDerivationKey()
	if err != nil {
		return nil, err
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
		return nil, err
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
		return err
	}
	defer db.Close()
	return db.DeleteLogin(l.Vault.Uuid, l.Name)
}

func (l *Login) Update(target string) (int64, error) {
	var cipher string
	var err error
	if l.Password != "" {
		e := crypto.New("aes")
		cipher, err = e.Encrypt(l.Vault.DerivationKey, l.Password)
		if err != nil {
			return 0, err
		}
	}
	db, err := storage.New("./database.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	aff, err := db.UpdateLogin(l.Vault.Uuid, target, l.Name, l.Username, cipher)
	if err != nil {
		return 0, err
	}
	return aff, nil
}
