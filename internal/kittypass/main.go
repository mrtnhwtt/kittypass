package kittypass

import (
	"fmt"

	"github.com/mrtnhwtt/kittypass/internal/crypto"
)

type KittyPass struct {
	Password      string
	DerivationKey []byte
}

func New() KittyPass {
	return KittyPass{}
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
