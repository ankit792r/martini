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

	if !needsInteractive(opts) {
		session := NewSession(opts.ProjectPath)
		if opts.IconsPath != "" {
			session.SetAnswers(CmdUpdateIcons, []string{opts.IconsPath})
		}
		return RunCommands(session, opts.Preselected)
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

func needsInteractive(opts Options) bool {
	if len(opts.Preselected) == 0 {
		return true
	}
	if opts.ProjectPath == "" {
		return true
	}
	for _, id := range opts.Preselected {
		cmd, ok := CommandByID(id)
		if !ok {
			continue
		}
		if id == CmdUpdateIcons && opts.IconsPath != "" {
			continue
		}
		if len(cmd.Fields) > 0 {
			return true
		}
	}
	return false
}
