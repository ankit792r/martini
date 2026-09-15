package devbox

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	stepProjectName step = iota
	stepAskEnv
	stepEnvKey
	stepEnvValue
	stepAskMoreEnv
	stepAskScript
	stepScriptName
	stepScriptCmd
	stepAskMoreScript
)

type mode int

const (
	modeInput mode = iota
	modeConfirm
)

type model struct {
	step        step
	mode        mode
	input       textinput.Model
	projectName string
	env         map[string]string
	scripts     map[string][]string
	pendingKey  string
	err         error
	quitting    bool
	done        bool
	outputPath  string
}

func newModel() model {
	return model{
		step:    stepProjectName,
		mode:    modeInput,
		env:     defaultEnv(),
		scripts: make(map[string][]string),
		input:   newInput("Project name: "),
	}
}

func defaultEnv() map[string]string {
	return map[string]string{
		"PATH": "$PATH:$HOME/.pub-cache/bin",
	}
}

func newInput(prompt string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.CharLimit = 256
	plain := lipgloss.NewStyle()
	ti.PromptStyle = plain
	ti.TextStyle = plain
	ti.Cursor.Style = plain
	ti.Focus()
	return ti
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		if m.mode == modeConfirm {
			next, cmd := m.handleConfirm(key)
			return next, cmd
		}

		if key.String() == "enter" {
			next := m.advance()
			if next.done {
				return next, tea.Quit
			}
			return next, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) handleConfirm(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m.confirmYes(), nil
	case "n", "N", "enter":
		next := m.confirmNo()
		if next.done {
			return next, tea.Quit
		}
		return next, nil
	}
	return m, nil
}

func (m model) advance() model {
	value := strings.TrimSpace(m.input.Value())

	switch m.step {
	case stepProjectName:
		if value == "" {
			return m
		}
		m.projectName = value
		m.step = stepAskEnv
		m.mode = modeConfirm
		m.input = textinput.Model{}

	case stepEnvKey:
		if value == "" {
			return m
		}
		m.pendingKey = value
		m.step = stepEnvValue
		m.input = newInput("Env value: ")

	case stepEnvValue:
		if value == "" {
			return m
		}
		m.env[m.pendingKey] = value
		m.pendingKey = ""
		m.step = stepAskMoreEnv
		m.mode = modeConfirm
		m.input = textinput.Model{}

	case stepScriptName:
		if value == "" {
			return m
		}
		m.pendingKey = value
		m.step = stepScriptCmd
		m.input = newInput("Script command: ")

	case stepScriptCmd:
		if value == "" {
			return m
		}
		m.scripts[m.pendingKey] = []string{value}
		m.pendingKey = ""
		m.step = stepAskMoreScript
		m.mode = modeConfirm
		m.input = textinput.Model{}
	}

	return m
}

func (m model) confirmYes() model {
	switch m.step {
	case stepAskEnv, stepAskMoreEnv:
		m.step = stepEnvKey
		m.mode = modeInput
		m.input = newInput("Env key: ")
	case stepAskScript, stepAskMoreScript:
		m.step = stepScriptName
		m.mode = modeInput
		m.input = newInput("Script name: ")
	}
	return m
}

func (m model) confirmNo() model {
	switch m.step {
	case stepAskEnv, stepAskMoreEnv:
		m.step = stepAskScript
		m.mode = modeConfirm
	case stepAskScript, stepAskMoreScript:
		return m.finish()
	}
	return m
}

func (m model) finish() model {
	cfg := buildConfig(m.projectName, m.env, m.scripts)
	path := "devbox.json"
	if err := writeConfig(path, cfg); err != nil {
		m.err = err
	} else {
		m.outputPath = path
	}
	m.done = true
	return m
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if m.done {
		if m.err != nil {
			return fmt.Sprintf("error: %v\n", m.err)
		}
		return fmt.Sprintf("Wrote %s\n", m.outputPath)
	}

	var b strings.Builder

	switch m.step {
	case stepProjectName:
		b.WriteString("Devbox setup\n\n")
		b.WriteString(m.input.View())
	case stepAskEnv:
		b.WriteString("Add custom environment variables? (y/n)\n")
	case stepEnvKey, stepEnvValue:
		b.WriteString(m.input.View())
	case stepAskMoreEnv:
		b.WriteString("Add another environment variable? (y/n)\n")
	case stepAskScript:
		b.WriteString("Add custom scripts? (y/n)\n")
	case stepScriptName, stepScriptCmd:
		b.WriteString(m.input.View())
	case stepAskMoreScript:
		b.WriteString("Add another script? (y/n)\n")
	}

	return b.String()
}
