package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/mrtnhwtt/kittypass/internal/kittypass"
)

type LoginsListModel struct {
	status status
	spinner spinner.Model
	list list.Model
	logins []kittypass.Login
}