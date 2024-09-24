package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
)

type LoginDisplayModel struct {
	status status
	spinner spinner.Model
}
