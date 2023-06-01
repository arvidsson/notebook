package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
    tabs []string
    activeTab int
    textarea textarea.Model
    err error
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
    border := lipgloss.NormalBorder()
    border.BottomLeft = left
    border.Bottom = middle
    border.BottomRight = right
    return border
}

var (
    inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
    activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
    docStyle          = lipgloss.NewStyle().Padding(0, 0, 0, 0)
    highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
    inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
    activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
    windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 1).Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func initialModel() model {
    text := textarea.New()
    text.Focus()
    tabs := []string{"Journal", "Inbox"}

    return model{
        tabs: tabs,
        textarea: text,
        err: nil,
    }
}

func (m *model) initModel(width int, height int) {
    // m.textarea.SetWidth(width)
    // m.textarea.SetHeight(height)
}

func (m model) Init() tea.Cmd {
    return tea.Batch(textarea.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.initModel(msg.Width, msg.Height)
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

func sum(str ...string) int {
    result := 0
    for _, s := range str {
        result += 1
        _ = s
    }
    return result
}

func (m model) View() string {
    doc := strings.Builder{}
    var renderedTabs []string

    var curWidth int

    for i, t := range m.tabs {
        var style lipgloss.Style
        isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
        if isActive {
            style = activeTabStyle.Copy()
        } else {
            style = inactiveTabStyle.Copy()
        }
        border, _, _, _, _ := style.GetBorder()
        if isFirst && isActive {
            border.BottomLeft = "│"
        } else if isFirst && !isActive {
            border.BottomLeft = "├"
        } else if isLast && isActive {
            border.BottomRight = "│"
        } else if isLast && !isActive {
            border.BottomRight = "┤"
        }
        style = style.Border(border)
        if i == (len(m.tabs)-1) {
            style.Width(80 - curWidth * 4)
        }
        renderedTabs = append(renderedTabs, style.Render(t))
        curWidth += sum(renderedTabs...)
    }

    row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
    doc.WriteString(row)
    doc.WriteString("\n")
    doc.WriteString(windowStyle.Width(80).Render(m.textarea.View()))
    return docStyle.Render(doc.String())
}

func main() {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())

    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}