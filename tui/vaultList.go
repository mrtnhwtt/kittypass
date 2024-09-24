package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
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
	status  status
	spinner spinner.Model
	list    list.Model
	vaults  []kittypass.Vault
}

func NewVaultList() VaultListModel {
	vl := VaultListModel{
		status: Loading,
	}
	vl.spinner = spinner.New()
	vl.spinner.Spinner = spinner.Line
	vl.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	vl.list = list.New(nil, list.NewDefaultDelegate(), 0,0)
	vl.list.Title = "Vaults"
	return vl
}

func (v VaultListModel) Init() tea.Cmd {
	return nil
}

func (v VaultListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case BlockSizeMsg:
		v.list.SetSize(msg.Width, msg.Height)
	}
	return nil, nil
}

func (v VaultListModel) View() string {
	return blockStyle.Render(v.list.View())
}