package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mrtnhwtt/kittypass/internal/kittypass"
)

var (
	blockStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, true)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type VaultListModel struct {
	status status
	list   list.Model
	vaults []kittypass.Vault
}

func NewVaultList() VaultListModel {
	vl := VaultListModel{
		status: Loading,
	}
	vl.list = list.New(nil, list.NewDefaultDelegate(), 0, 0)
	vl.list.Title = "Vaults"
	return vl
}

func (v VaultListModel) Init() tea.Cmd {
	return nil
}

func (v VaultListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case BlockSizeMsg:
		// v.list.SetSize(msg.Width, msg.Height)
		v.list.SetWidth(msg.Width)
		v.list.SetHeight(msg.Height)
		blockStyle = blockStyle.Width(msg.Width)
		blockStyle = blockStyle.Height(msg.Height)
	case LoadedVaultListMsg:
		v.vaults = msg.vaults
		var newItems []list.Item
		for _, vault := range(v.vaults) {
			newItems = append(newItems, item{title: vault.Name, description: vault.Description})
		}
		v.status = Visible
		return v, v.list.SetItems(newItems)
	}
	return v, nil
}

func (v VaultListModel) View() string {
	if v.status == Visible {
		return blockStyle.Render(v.list.View())
	}
	return "invalid status"
}

