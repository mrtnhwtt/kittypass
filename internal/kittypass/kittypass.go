package kittypass

import (
	"encoding/hex"
	"fmt"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
	"github.com/mrtnhwtt/kittypass/internal/storage"
)

type KittyPass struct {
	Password      string
	Username      string
	Name          string
	DerivationKey []byte
	Salt          []byte
}

func New() KittyPass {
	return KittyPass{}
}

// UseMasterPassword generate a DerivationKey from a master password to use to encrypt and decrypt passwords stored in kittypass
func (kp *KittyPass) UseMasterPassword(masterPassword string) {
	kp.Salt = crypto.GenerateRandomSalt(16)
	kp.DerivationKey = crypto.GenerateKey([]byte(masterPassword), kp.Salt)
}

func (kp *KittyPass) RecreateDerivationKey(saltHex, masterPassword string) error {
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return fmt.Errorf("error when decoding salt")
	}
	kp.Salt = salt
	kp.DerivationKey = crypto.GenerateKey([]byte(masterPassword), kp.Salt)
	return nil
}

func (kp *KittyPass) Add() error {
	e := crypto.New("aes")
	cipher, err := e.Encrypt(kp.DerivationKey, kp.Password)
	if err != nil {
		return fmt.Errorf("error while encrypting password: %s", err)
	}
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}
	hexSalt := hex.EncodeToString(kp.Salt)
	_, err = db.Save(kp.Name, kp.Username, cipher, hexSalt)
	if err != nil {
		return fmt.Errorf("error while saving to database: %s", err)
	}
	return nil
}

func (kp *KittyPass) Get(master string) (map[string]string, error) {
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	stored, err := db.Read(kp.Name)
	if err != nil {
		return nil, fmt.Errorf("error reading stored password: %s", err)
	}
	saltHex := stored["salt"]
	err = kp.RecreateDerivationKey(saltHex, master)
	if err != nil {
		return nil, fmt.Errorf("error while creating derivation key from stored salt and masterpassword: %s", err)
	}
	e := crypto.New("aes")
	decrypted, err := e.Decrypt(kp.DerivationKey, stored["encrypted_password"])
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

func (kp *KittyPass) List() ([]map[string]string, error) {
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	loginList, err := db.List(kp.Name, kp.Username)
	if err != nil {
		return nil, err
	}
	return loginList, nil
}

func (kp *KittyPass) Run(input string) error {
	fmt.Println(input)
	salt := crypto.GenerateRandomSalt(16)
	kp.DerivationKey = crypto.GenerateKey([]byte(kp.Password), salt)

	e := crypto.New("aes")
	cipher, err := e.Encrypt(kp.DerivationKey, input)
	if err != nil {
		return err
	}
	fmt.Println(cipher)
	decrypted, err := e.Decrypt(kp.DerivationKey, string(cipher))
	if err != nil {
		return err
	}
	fmt.Println(decrypted)
	return nil
}
