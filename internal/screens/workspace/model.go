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

func New() Model {
	testCases, _ := core.GetTestCases()
	testCaseList := list.New(&testCases)

	model := Model{
		TestCaseList: testCaseList,
		Input:        textarea.New("Input"),
		Expected:     textarea.New("Expected Output"),
		Output:       output.New(),
		Help:         help.New(),

		keyMap: DefaultKeyMap(),
	}

	model.focusOn(0)

	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor, textarea.Blink, m.TestCaseList.SelectTestCase(0))
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
			// help height is dynamic
			m.updateLayout()
		}
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

	cmds = append(cmds, m.currentlyFocused().Update(msg))

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	help := m.HelpView()

	credit := lipgloss.NewStyle().Italic(true).Hyperlink("https://thecomputerm.dev").Render("by TheComputerM")

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
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			help,
			lipgloss.PlaceHorizontal(m.width-lipgloss.Width(help), lipgloss.Right, credit),
		),
	)

	return zone.Scan(content)
}
