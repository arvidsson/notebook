package main

import (
    "fmt"
    "os"

    "github.com/charmbracelet/bubbles/textarea"
    tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
    textarea textarea.Model
    err error
}

func initialModel() model {
    text := textarea.New()
    text.Placeholder = "Start writing text here..."
    text.Focus()

	return model{
        textarea: text,
        err: nil,
	}
}

func (m model) Init() tea.Cmd {
    return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
    return fmt.Sprintf(
        "%s",
        m.textarea.View(),
    ) + "\n"
}

func main() {
    p := tea.NewProgram(initialModel())

    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}