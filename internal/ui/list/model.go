package list

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/core"
)

type Model struct {
	items  []core.TestCase
	keyMap KeyMap
	styles Styles

	index   int
	focused bool
}

func New(testCases []core.TestCase) Model {
	return Model{
		items:  testCases,
		keyMap: DefaultKeyMap(),
		styles: DefaultStyles(true),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.focused {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keyMap.Up):
			m.index = max(m.index-1, 0)
		case key.Matches(msg, m.keyMap.Down):
			m.index = min(m.index+1, len(m.items)-1)
		case key.Matches(msg, m.keyMap.Create):
			core.CreateTestCase(&m.items)
			m.index = len(m.items) - 1
		}
	}
	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(m.styles.Title.Render("TITLE"))
	sb.WriteByte('\n')

	var state StyleState = m.styles.Blurred
	if m.focused {
		state = m.styles.Focused
	}

	for i, item := range m.items {
		if i > 0 {
			sb.WriteString("\n\n")
		}

		content := fmt.Sprintf(
			"%s\n%s",
			m.styles.getTitleStyle(item.Status).Bold(i == m.index).Render(string(item.Status)),
			state.ItemDesc.Faint(i != m.index).Render(item.Details),
		)

		if i != m.index {
			// normal item
			sb.WriteString(m.styles.Item.Render(content))
		} else {
			// selected item
			sb.WriteString(state.SelectedItem.Render(content))
		}

	}

	return m.styles.List.Render(sb.String())
}

func (m *Model) SetStyles(styles Styles) {
	m.styles = styles
}

func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *Model) Blur() {
	m.focused = false
}
