package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	status     status
	spinner    spinner.Model
	vaultList  VaultListModel
	loginsList LoginsListModel
	login      LoginDisplayModel
	db         *storage.Storage
	conf       *viper.Viper
}

func NewBoard(v *viper.Viper) *Board {
	vl := NewVaultList()
	board := &Board{
		conf:      v,
		vaultList: vl,
		status: Loading,
	}
	board.spinner = spinner.New()
	board.spinner.Spinner = spinner.Line
	board.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return board
}

func (b Board) Init() tea.Cmd {
	return tea.Batch(loadDb(b.conf))
}

func (b Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		b.vaultList.Update(BlockSizeMsg{Width: msg.Width / 3, Height: msg.Height})
	case LoadDatabaseErrorMsg:
		b.focus = ErrorView
	case loadedDatabaseMsg:
		b.db = msg.db
		return b, LoadVaultList(b.db)
	case LoadVaultsErrorMsg:
		b.focus = ErrorView
	case LoadedVaultListMsg:
		var cmd tea.Cmd
		var m tea.Model
		m, cmd = b.vaultList.Update(msg)
		b.status = Visible
		if vm, ok := m.(VaultListModel); ok {
			b.vaultList = vm
		}
		return b, cmd
	}

	var cmd tea.Cmd
	b.spinner, cmd = b.spinner.Update(msg)
	return b, cmd
}

func (b Board) View() string {
	if b.status == Visible {
		return lipgloss.JoinHorizontal(lipgloss.Center, b.vaultList.View())
	}
	return b.spinner.View()
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
		for _, vaultInfo := range vaultsMap {
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
