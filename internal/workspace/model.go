package workspace

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	zone "github.com/lrstanley/bubblezone/v2"
	"github.com/thecomputerm/lazycph/internal/core"
	"github.com/thecomputerm/lazycph/internal/ui/list"
	"github.com/thecomputerm/lazycph/internal/ui/output"
	"github.com/thecomputerm/lazycph/internal/ui/textarea"
)

type Model struct {
	TestCaseList list.Model
	Input        textarea.Model
	Expected     textarea.Model
	Output       output.Model
	Help         help.Model

	keyMap KeyMap

	focused       uint
	width, height int
}

var _ tea.Model = (*Model)(nil)

func New() Model {
	testCases, _ := core.GetTestCases()
	testCaseList := list.New(testCases)

	model := Model{
		TestCaseList: testCaseList,
		Input:        textarea.New("Input"),
		Expected:     textarea.New("Expected Output"),
		Output:       output.New(),
		Help:         help.New(),
		keyMap:       DefaultKeyMap(),
	}

	model.focusOn(0)

	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor, textarea.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.setTheme(msg.IsDark())
		return m, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Next):
			return m, m.focusNext()
		case key.Matches(msg, m.keyMap.Prev):
			return m, m.focusPrev()
		case key.Matches(msg, m.keyMap.Help):
			m.Help.ShowAll = !m.Help.ShowAll
			m.updateLayout()
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.MouseReleaseMsg:
		if msg.Button == tea.MouseLeft {
			if zone.Get("section-list").InBounds(msg) {
				cmds = append(cmds, m.focusOn(0))
			} else if zone.Get("section-input").InBounds(msg) {
				cmds = append(cmds, m.focusOn(1))
			} else if zone.Get("section-expected").InBounds(msg) {
				cmds = append(cmds, m.focusOn(2))
			} else if zone.Get("section-output").InBounds(msg) {
				cmds = append(cmds, m.focusOn(3))
			}
		}
	}

	m.TestCaseList, cmd = m.TestCaseList.Update(msg)
	cmds = append(cmds, cmd)

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	m.Expected, cmd = m.Expected.Update(msg)
	cmds = append(cmds, cmd)

	m.Output, cmd = m.Output.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			zone.Mark("section-list", m.TestCaseList.View()),
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					zone.Mark("section-input", m.Input.View()),
					zone.Mark("section-expected", m.Expected.View()),
				),
				zone.Mark("section-output", m.Output.View()),
			),
		),
		m.HelpView(),
	)

	v := tea.NewView(zone.Scan(content))

	v.AltScreen = true
	v.WindowTitle = "LazyCPH"
	v.MouseMode = tea.MouseModeCellMotion

	return v
}

func (m *Model) setTheme(isDark bool) {
	m.TestCaseList.SetStyles(list.DefaultStyles(isDark))
	m.Input.SetStyles(textarea.DefaultStyles(isDark))
	m.Expected.SetStyles(textarea.DefaultStyles(isDark))
	m.Output.SetStyles(output.DefaultStyles(isDark))
}
