package flutter

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func newInput(prompt string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.CharLimit = 512
	plain := lipgloss.NewStyle()
	ti.PromptStyle = plain
	ti.TextStyle = plain
	ti.Cursor.Style = plain
	ti.Focus()
	return ti
}
