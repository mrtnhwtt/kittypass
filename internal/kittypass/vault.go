package kittypass

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
	"github.com/mrtnhwtt/kittypass/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Vault struct {
	Uuid              string
	Name              string
	Description       string
	Masterpass        string
	HexHashMasterpass string
	HexSalt           string
	DerivationKey     []byte
	Salt              []byte
}

func NewVault() Vault {
	return Vault{}
}

// UseMasterPassword generate a DerivationKey from a master password to use to encrypt and decrypt passwords stored in kittypass
func (v *Vault) UseMasterPassword() {
	v.Salt = crypto.GenerateRandomSalt(16)
	v.HexSalt = hex.EncodeToString(v.Salt)
	v.DerivationKey = crypto.GenerateKey([]byte(v.Masterpass), v.Salt)
}

func (v *Vault) RecreateDerivationKey() error {
	salt, err := hex.DecodeString(v.HexSalt)
	if err != nil {
		return fmt.Errorf("error when decoding salt")
	}
	v.Salt = salt
	v.DerivationKey = crypto.GenerateKey([]byte(v.Masterpass), v.Salt)
	return nil
}

func (v *Vault) CreateVault() error {
	v.UseMasterPassword()
	err := v.HashMasterpass()
	if err != nil {
		return err
	}
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}
	defer db.Close()

	_, err = db.SaveVault(v.Name, v.Description, v.HexHashMasterpass, v.HexSalt)
	if err != nil {
		return errors.New("error while generating vault, please try again")
	}
	return nil
}

func (v *Vault) HashMasterpass() error {
	hashedMaster, err := bcrypt.GenerateFromPassword([]byte(v.Masterpass), 15)
	if err != nil {
		return err
	}
	v.HexHashMasterpass = hex.EncodeToString(hashedMaster)
	return nil
}

// Check the user provided master password against the one defined when creating the vault.
// If they match, user can interact with logins in the vault. Returns an error when they do not match and nil when they match
func (v *Vault) MasterpassMatch() error {
	storedPass, err := hex.DecodeString(v.HexHashMasterpass)
	if err != nil {
		return err
	}
	// check the stored master password against the provided master password
	if err = bcrypt.CompareHashAndPassword(storedPass, []byte(v.Masterpass)); err != nil {
		return err
	}
	return nil
}

func (v *Vault) Get() error {
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}
	defer db.Close()
	vaultData, err := db.GetVault(v.Name)
	if err != nil {
		return err
	}
	v.Uuid = vaultData["uuid"]
	v.Description = vaultData["description"]
	v.HexHashMasterpass = vaultData["hex_hashed_master_password"]
	v.HexSalt = vaultData["hex_salt"]
	v.Salt, err = hex.DecodeString(v.HexSalt)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) List() ([]map[string]string, error) {
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	defer db.Close()
	vaultList, err := db.ListVault(v.Name)
	if err != nil {
		return nil, err
	}
	return vaultList, nil
}

func (v *Vault) Delete() (map[string]int64, error) {
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	defer db.Close()
	deleted, err := db.DeleteVault(v.Name, v.Uuid)
	if err != nil {
		return nil, err
	}
	return deleted, nil
}

func (v *Vault) Update(newMasterPass, newName, newDescription string) (map[string]int, error) {
	var err error
	var loginList []map[string]string
	db, err := storage.New("./database.db")
	if err != nil {
		return nil, fmt.Errorf("error while connecting to storage: %s", err)
	}
	defer db.Close()
	if newMasterPass != "" {
		loginList, err = v.reencryptLogins(db, newMasterPass)
		if err != nil {
			return nil, err
		}

	}
	return db.UpdateVault(v.Uuid, newName, newDescription, v.HexHashMasterpass, v.HexSalt, loginList)

}

func (v *Vault) reencryptLogins(db *storage.Storage, newMasterPass string) ([]map[string]string, error) {
	// get all login and decrypt the passwords
	e := crypto.New("aes")
	loginList, err := db.ReadLogins(v.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get login for vault password update: %s", err)
	}
	for _, login := range loginList {
		login["decrypted"], err = e.Decrypt(v.DerivationKey, login["hex_enc_pass"])
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt password: %s", err)
		}
	}

	// create a new derivation key from the new password and encrypt login passwords
	v.Masterpass = newMasterPass
	v.UseMasterPassword()
	err = v.HashMasterpass()
	if err != nil {
		return nil, err
	}
	for _, login := range loginList {
		login["newHexEncrypted"], err = e.Encrypt(v.DerivationKey, login["decrypted"])
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt password: %s", err)
		}
		delete(login, "decrypted")
	}
	return loginList, nil
}
