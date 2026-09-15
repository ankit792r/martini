package flutter

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(opts Options) error {
	if len(opts.Preselected) == 0 && opts.ProjectPath != "" {
		return fmt.Errorf("project path requires a command flag (e.g. --update-package-name)")
	}

	p := tea.NewProgram(
		newModel(opts),
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stdout),
	)
	final, err := p.Run()
	if err != nil {
		return err
	}

	m, ok := final.(*model)
	if !ok {
		return nil
	}
	return m.err
}
