package devbox

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
	if _, err := os.Stat("devbox.json"); err == nil {
		return fmt.Errorf("devbox.json already exists in current directory")
	}

	p := tea.NewProgram(
		newModel(),
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stdout),
	)
	final, err := p.Run()
	if err != nil {
		return err
	}

	m, ok := final.(model)
	if !ok {
		return nil
	}
	return m.err
}
