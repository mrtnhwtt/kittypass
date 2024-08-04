package kittypass

import (
	"encoding/hex"
	"fmt"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
	"github.com/mrtnhwtt/kittypass/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Vault struct {
	Uuid          string
	Name          string
	Description   string
	Masterpass    string
	HexHash       string
	DerivationKey []byte
	Salt          []byte
}

func NewVault() Vault {
	return Vault{}
}

// UseMasterPassword generate a DerivationKey from a master password to use to encrypt and decrypt passwords stored in kittypass
func (v *Vault) UseMasterPassword() {
	v.Salt = crypto.GenerateRandomSalt(16)
	v.DerivationKey = crypto.GenerateKey([]byte(v.Masterpass), v.Salt)
}

func (v *Vault) RecreateDerivationKey(saltHex, masterPassword string) error {
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return fmt.Errorf("error when decoding salt")
	}
	v.Salt = salt
	v.DerivationKey = crypto.GenerateKey([]byte(masterPassword), v.Salt)
	return nil
}

func (v *Vault) CreateVault() error {
	err := v.HashMasterpass()
	if err != nil {
		return err
	}
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}
	hexSalt := hex.EncodeToString(v.Salt)
	db.NewVault(v.Name, v.Description, v.HexHash, hexSalt)
	return nil
}

func (v *Vault) HashMasterpass() error {
	hashedMaster, err := bcrypt.GenerateFromPassword([]byte(v.Masterpass), 15)
	if err != nil {
		return err
	}
	v.HexHash = hex.EncodeToString(hashedMaster)
	return nil
}

// Check the user provided master password against the one defined when creating the vault.
// If they match, user can interact with logins in the vault. Returns an error when they do not match and nil when they match
func (v *Vault) MasterpassMatch() error {
	db, err := storage.New("./database.db")
	if err != nil {
		return fmt.Errorf("error while connecting to storage: %s", err)
	}
	vaultData, err := db.GetVault(v.Name)
	if err != nil {
		return err
	}
	storedPass, err := hex.DecodeString(vaultData["hashed_master_password"])
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
	vaultData, err := db.GetVault(v.Name)
	if err != nil {
		return err
	}
	v.Uuid = vaultData["uuid"]
	v.Description = vaultData["description"]
	v.HexHash = vaultData["hashed_master_password"]
	v.Salt, err = hex.DecodeString(vaultData["salt"])
	if err != nil {
		return err
	}
	return nil
}