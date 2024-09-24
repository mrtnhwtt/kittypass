package kittypass

import (
	"encoding/hex"
	"log"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
	"github.com/mrtnhwtt/kittypass/internal/storage"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type Vault struct {
	Uuid              string
	Name              string
	Description       string
	Masterpass        string
	HexHashMasterpass string
	Db                *storage.Storage
	DateCreated       string
}

func NewVault() Vault {
	return Vault{}
}

func (v *Vault) OpenDbConnection(conf *viper.Viper) error {
	var err error
	v.Db, err = storage.New(conf.GetString("storage_path"))
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) CreateVault() error {
	err := v.HashMasterpass()
	if err != nil {
		return err
	}

	_, err = v.Db.SaveVault(v.Name, v.Description, v.HexHashMasterpass)
	if err != nil {
		return err
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
		log.Printf("error when decoding stored hex cipher: %s", err)
		return MalformedDataError{"hex cipher"}
	}
	// check the stored master password against the provided master password
	if err = bcrypt.CompareHashAndPassword(storedPass, []byte(v.Masterpass)); err != nil {
		return IncorrectPasswordError{}
	}
	return nil
}

func (v *Vault) Get() error {
	vaultData, err := v.Db.GetVault(v.Name)
	if err != nil {
		return err
	}
	v.Uuid = vaultData["uuid"]
	v.Description = vaultData["description"]
	v.HexHashMasterpass = vaultData["hex_cipher"]
	return nil
}

func (v *Vault) List() ([]map[string]string, error) {
	vaultList, err := v.Db.ListVault(v.Name)
	if err != nil {
		return nil, err
	}
	return vaultList, nil
}

func (v *Vault) Delete() (map[string]int64, error) {
	deleted, err := v.Db.DeleteVault(v.Name, v.Uuid)
	if err != nil {
		return nil, err
	}
	return deleted, nil
}

func (v *Vault) Update(newMasterPass, newName, newDescription string) (map[string]int, error) {
	var err error
	var loginList []map[string]string
	defer v.Db.Close()
	if newMasterPass != "" {
		loginList, err = v.reencryptLogins(newMasterPass)
		if err != nil {
			return nil, err
		}

	}
	return v.Db.UpdateVault(v.Uuid, newName, newDescription, v.HexHashMasterpass, loginList)
}

func (v *Vault) reencryptLogins(newMasterPass string) ([]map[string]string, error) {
	// get all login and decrypt the passwords
	e := crypto.New("aes")
	loginList, err := v.Db.ReadLogins(v.Uuid)
	if err != nil {
		return nil, err
	}
	var login Login
	for _, loginData := range loginList {
		login = NewLogin()
		login.HexSalt = loginData["hex_salt"]
		login.Vault = v
		login.RecreateDerivationKey()
		loginData["decrypted"], err = e.Decrypt(login.DerivationKey, loginData["hex_cipher"])
		if err != nil {
			return nil, err
		}
	}

	// create a new derivation key from the new password and encrypt login passwords
	v.Masterpass = newMasterPass
	err = v.HashMasterpass()
	if err != nil {
		return nil, err
	}
	for _, loginData := range loginList {
		login = NewLogin()
		login.Vault = v
		login.UseMasterPassword()
		loginData["new_hex_cipher"], err = e.Encrypt(login.DerivationKey, loginData["decrypted"])
		if err != nil {
			return nil, err
		}
		loginData["new_hex_salt"] = login.HexSalt
		delete(loginData, "decrypted")
	}
	return loginList, nil
}
