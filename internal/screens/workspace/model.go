package workspace

import (
	"path/filepath"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	zone "github.com/lrstanley/bubblezone/v2"
	"github.com/thecomputerm/lazycph/internal/core"
	"github.com/thecomputerm/lazycph/internal/screens/filepicker"
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

	filePath      string
	focused       uint
	width, height int
}

var _ tea.Model = (*Model)(nil)

func New(msg filepicker.FileSelectedMsg) Model {
	zone.NewGlobal()

	model := Model{
		TestCaseList: list.New(),
		Input:        textarea.New("Input"),
		Expected:     textarea.New("Expected Output"),
		Output:       output.New(),
		Help:         help.New(),

		keyMap:   DefaultKeyMap(),
		filePath: msg.Path,
	}

	model.TestCaseList.Title = filepath.Base(msg.Path)
	model.TestCaseList.Items = msg.TestCases

	model.focusOn(0)

	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
		tea.RequestBackgroundColor,
		textarea.Blink,
		m.TestCaseList.SelectTestCase(0),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// set app theme
		isDark := msg.IsDark()
		m.TestCaseList.SetStyles(list.DefaultStyles(isDark))
		m.Input.SetStyles(textarea.DefaultStyles(isDark))
		m.Expected.SetStyles(textarea.DefaultStyles(isDark))
		m.Output.SetStyles(output.DefaultStyles(isDark))
		return m, nil
	case tea.WindowSizeMsg:
		// set window size
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		return m, nil
	case list.TestCaseSelectedMsg:
		// new testcase selected
		m.Input.BindValue(&msg.TestCase.Input)
		m.Expected.BindValue(&msg.TestCase.Expected)
		m.Output.SetContent(msg.TestCase.Output)
	case core.TestCaseExecutedMsg:
		// testcase finished running; refresh output if still selected
		if m.TestCaseList.Selected() == msg.TestCase {
			m.Output.SetContent(msg.TestCase.Output)
		}

		// save updated testcase output
		return m, m.TestCaseList.Items.Save(m.filePath)
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
		case key.Matches(msg, m.keyMap.Run):
			return m, m.TestCaseList.Selected().Execute(m.filePath)
		case key.Matches(msg, m.keyMap.RunAll):
			for _, tc := range m.TestCaseList.Items {
				cmds = append(cmds, tc.Execute(m.filePath))
			}
			return m, tea.Batch(cmds...)
		}
	case tea.MouseReleaseMsg:
		if msg.Button == tea.MouseLeft {
			if zone.Get("section-list").InBounds(msg) {
				cmds = append(cmds, m.focusOn(0))
			} else if zone.Get("section-input").InBounds(msg) {
				cmds = append(cmds, m.focusOn(1))
			} else if zone.Get("section-expected").InBounds(msg) {
				cmds = append(cmds, m.focusOn(2))
			} else if zone.Get("setion-output").InBounds(msg) {
				cmds = append(cmds, m.focusOn(3))
			} else if zone.Get("section-help").InBounds(msg) {
				m.Help.ShowAll = !m.Help.ShowAll
				m.updateLayout()
			}
		}
	}

	var cmd tea.Cmd
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
		zone.Mark("section-help", m.helpView()),
	)

	v := tea.NewView(zone.Scan(content))

	v.AltScreen = true
	v.WindowTitle = "LazyCPH"
	v.MouseMode = tea.MouseModeCellMotion

	return v
}
