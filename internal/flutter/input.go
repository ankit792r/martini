package flutter

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func newInput(prompt string) textinput.Model {
	return newInputWithOptions(prompt, false)
}

func newSecretInput(prompt string) textinput.Model {
	return newInputWithOptions(prompt, true)
}

func newInputWithOptions(prompt string, secret bool) textinput.Model {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.CharLimit = 512
	if secret {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '*'
	}
	plain := lipgloss.NewStyle()
	ti.PromptStyle = plain
	ti.TextStyle = plain
	ti.Cursor.Style = plain
	ti.Focus()
	return ti
}
