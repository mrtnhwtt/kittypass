package tui

import (
	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/storage"
)

type loadedDatabaseMsg struct {
	db *storage.Storage
}
type LoadDatabaseErrorMsg struct {
	err error
}
func (e LoadDatabaseErrorMsg) Error() string { return e.err.Error() }

type LoadedVaultListMsg struct {
	vaults []kittypass.Vault
}

type LoadVaultsErrorMsg struct {
	err error
}

func (e LoadVaultsErrorMsg) Error() string { return e.err.Error() }

type BlockSizeMsg struct {
	Width int
	Height int
}