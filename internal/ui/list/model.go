package list

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
	"github.com/thecomputerm/lazycph/internal/core"
)

type Model struct {
	Title string

	Items  core.TestCaseList
	KeyMap KeyMap
	styles Styles

	// index of the currently selected test case
	index   int
	focused bool
	height  int
}

func New(title string) Model {
	return Model{
		Title: title,

		KeyMap: DefaultKeyMap(),
		styles: DefaultStyles(true),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Up):
			return m, m.SelectTestCase(max(m.index-1, 0))
		case key.Matches(msg, m.KeyMap.Down):
			return m, m.SelectTestCase(min(m.index+1, len(m.Items)-1))
		case key.Matches(msg, m.KeyMap.Create):
			m.Items.Create()
			m.KeyMap.Delete.SetEnabled(true)
			return m, m.SelectTestCase(len(m.Items) - 1)
		case key.Matches(msg, m.KeyMap.Delete):
			m.Items.Delete(m.index)
			if len(m.Items) == 1 {
				// Disable delete key when only one test case remains
				m.KeyMap.Delete.SetEnabled(false)
			}
			return m, m.SelectTestCase(max(m.index-1, 0))
		}
	case tea.MouseReleaseMsg:
		if msg.Button == tea.MouseLeft {
			for i, _ := range m.Items {
				if zone.Get("listitem-" + strconv.Itoa(i)).InBounds(msg) {
					return m, m.SelectTestCase(i)
				}
			}
		}
	}

	return m, nil
}

// Returns the currently selected test case
func (m *Model) Selected() *core.TestCase {
	return m.Items[m.index]
}

type TestCaseSelectedMsg struct {
	Index    int
	TestCase *core.TestCase
}

// Selects the test case at the given index and returns a TestCaseSelectedMsg
func (m *Model) SelectTestCase(index int) tea.Cmd {
	if !(index >= 0 && index < len(m.Items)) {
		return nil
	}

	m.index = index
	return func() tea.Msg {
		return TestCaseSelectedMsg{Index: index, TestCase: m.Selected()}
	}
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(m.styles.Title.Render(m.Title))
	sb.WriteString("\n\n")

	var state StyleState = m.styles.Blurred
	if m.focused {
		state = m.styles.Focused
	}

	for i, testCase := range m.Items {
		if i > 0 {
			sb.WriteString("\n\n")
		}

		content := fmt.Sprintf(
			"%s\n%s",
			m.styles.itemTitleStyle(testCase.Status).Bold(i == m.index).Render(string(testCase.Status)),
			state.ItemDesc.Faint(i != m.index).Render(testCase.Details),
		)

		var item string

		if i != m.index {
			// normal item
			item = m.styles.Item.Render(content)
		} else {
			// selected item
			item = state.SelectedItem.Render(content)
		}

		sb.WriteString(zone.Mark("listitem-"+strconv.Itoa(i), item))
	}

	return m.styles.List.Height(m.height).Render(sb.String())
}

func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *Model) Blur() {
	m.focused = false
}
