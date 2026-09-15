package flutter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type phase int

const (
	phaseSelect phase = iota
	phaseProjectPath
	phaseDetails
	phaseExecute
	phaseDone
)

type model struct {
	phase           phase
	cursor          int
	selected        map[CommandID]bool
	selectedOrder   []CommandID
	projectPath     string
	skipProjectPath bool
	input           textinput.Model
	commandIndex    int
	promptIndex     int
	session         *Session
	err             error
	quitting        bool
	done            bool
}

type Options struct {
	ProjectPath string
	Preselected []CommandID
}

func newModel(opts Options) *model {
	selected := make(map[CommandID]bool)
	order := append([]CommandID(nil), opts.Preselected...)

	for _, id := range order {
		selected[id] = true
	}

	m := &model{
		selected:        selected,
		selectedOrder:   order,
		projectPath:     opts.ProjectPath,
		skipProjectPath: opts.ProjectPath != "",
		session:         NewSession(opts.ProjectPath),
	}

	if len(order) == 0 {
		m.phase = phaseSelect
		return m
	}

	if m.skipProjectPath {
		m.phase = phaseDetails
		m.startNextPrompt()
		return m
	}

	m.phase = phaseProjectPath
	m.input = newInput("Flutter project path: ")
	return m
}

func (m *model) startNextPrompt() {
	for m.commandIndex < len(m.selectedOrder) {
		cmd, ok := CommandByID(m.selectedOrder[m.commandIndex])
		if !ok {
			m.commandIndex++
			m.promptIndex = 0
			continue
		}

		if m.promptIndex < len(cmd.Fields) {
			field := cmd.Fields[m.promptIndex]
			if field.Secret {
				m.input = newSecretInput(field.Label)
			} else {
				m.input = newInput(field.Label)
			}
			return
		}

		m.commandIndex++
		m.promptIndex = 0
	}

	m.phase = phaseExecute
}

func (m *model) Init() tea.Cmd {
	if m.phase == phaseSelect || m.phase == phaseExecute {
		return nil
	}
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if key.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.phase {
		case phaseSelect:
			return m.updateSelect(key)
		case phaseProjectPath, phaseDetails:
			if key.String() == "enter" {
				if !m.advanceInput() {
					return m, nil
				}
				if m.phase == phaseExecute {
					if err := RunCommands(m.session, m.selectedOrder); err != nil {
						m.err = err
					}
					m.phase = phaseDone
					m.done = true
					return m, tea.Quit
				}
				return m, textinput.Blink
			}
		}
	}

	if m.phase != phaseProjectPath && m.phase != phaseDetails {
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) updateSelect(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(AllCommands)-1 {
			m.cursor++
		}
	case "tab":
		id := AllCommands[m.cursor].ID
		if m.selected[id] {
			delete(m.selected, id)
			m.selectedOrder = removeID(m.selectedOrder, id)
		} else {
			m.selected[id] = true
			m.selectedOrder = append(m.selectedOrder, id)
		}
	case "enter":
		if len(m.selectedOrder) == 0 {
			return m, nil
		}
		if m.skipProjectPath {
			m.session.ProjectPath = m.projectPath
			m.phase = phaseDetails
			m.startNextPrompt()
			return m, textinput.Blink
		}
		m.phase = phaseProjectPath
		m.input = newInput("Flutter project path: ")
		return m, textinput.Blink
	}

	return m, nil
}

func (m *model) advanceInput() bool {
	value := strings.TrimSpace(m.input.Value())
	if value == "" {
		return false
	}

	switch m.phase {
	case phaseProjectPath:
		m.projectPath = value
		m.session.ProjectPath = value
		m.phase = phaseDetails
		m.startNextPrompt()

	case phaseDetails:
		id := m.selectedOrder[m.commandIndex]
		cmd, ok := CommandByID(id)
		if !ok {
			return false
		}

		field := cmd.Fields[m.promptIndex]
		if value == "" {
			if field.Required {
				return false
			}
			if field.Default != "" {
				value = field.Default
			}
		}

		answers := m.session.AnswersFor(id)
		answers = append(answers, value)
		m.session.SetAnswers(id, answers)
		m.promptIndex++
		m.startNextPrompt()
	}

	return true
}

func removeID(ids []CommandID, target CommandID) []CommandID {
	out := ids[:0]
	for _, id := range ids {
		if id != target {
			out = append(out, id)
		}
	}
	return out
}

func (m *model) View() string {
	if m.quitting {
		return ""
	}
	if m.done {
		if m.err != nil {
			return fmt.Sprintf("error: %v\n", m.err)
		}
		return "Done\n"
	}

	var b strings.Builder

	switch m.phase {
	case phaseSelect:
		b.WriteString("Flutter tasks\n")
		b.WriteString("tab: toggle  enter: continue  up/down: move\n\n")
		for i, cmd := range AllCommands {
			prefix := "  "
			if i == m.cursor {
				prefix = "> "
			}
			mark := "[ ]"
			if m.selected[cmd.ID] {
				mark = "[x]"
			}
			fmt.Fprintf(&b, "%s%s %s\n", prefix, mark, cmd.Label)
		}
	case phaseProjectPath, phaseDetails:
		b.WriteString(m.input.View())
	}

	return b.String()
}
