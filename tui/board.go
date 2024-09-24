package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/storage"
	"github.com/spf13/viper"
)

type focus int

const (
	VaultListView focus = iota
	LoginsListView
	LoginDetailView
	ErrorView
)

type status int

const (
	Loading status = iota
	Locked
	Visible
)

type Board struct {
	focus      focus
	vaultList  VaultListModel
	loginsList LoginsListModel
	login      LoginDisplayModel
	db         *storage.Storage
	conf       *viper.Viper
}

func NewBoard(v *viper.Viper) *Board {
	board := &Board{
		conf: v,
	}
	return board
}

func (b Board) Init() tea.Cmd {
	return tea.Batch(loadDb(b.conf))
}

func (b Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		// h, v := docStyle.GetFrameSize()
		// b.list.SetSize(msg.Width-h, msg.Height-v)
	case LoadDatabaseErrorMsg:
		b.focus = ErrorView
	case loadedDatabaseMsg:
		b.db = msg.db
	case LoadVaultsErrorMsg:
		b.focus = ErrorView
	case LoadedVaultListMsg:
		b.vaultList.vaults = msg.vaults
	}

	var cmd tea.Cmd
	// b, cmd = b.Update(msg)
	return b, cmd
}

func (b Board) View() string {
	return ""
}

func loadDb(v *viper.Viper) tea.Cmd {
    return func() tea.Msg {
        db, err := storage.New(v.GetString("storage_path"))
        if err != nil {
            return LoadDatabaseErrorMsg{err}
        }
        return loadedDatabaseMsg{db: db}
    }
}

func LoadVaultList(db *storage.Storage) tea.Cmd {
	return func() tea.Msg {
		vaultsMap, err := db.ListVault("")
		if err != nil {
			return LoadVaultsErrorMsg{err}
		}
		var vaults []kittypass.Vault
		for _, vaultInfo := range(vaultsMap) {
			vault := kittypass.NewVault()
			vault.Db = db
			vault.Description = vaultInfo["description"]
			vault.Name = vaultInfo["name"]
			vault.DateCreated = vaultInfo["date_created"]
			vaults = append(vaults, vault)
		}
		return LoadedVaultListMsg{vaults}
	}
}