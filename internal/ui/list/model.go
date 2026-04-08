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
	items  []core.TestCase
	keyMap KeyMap
	styles Styles

	index   int
	focused bool
	height  int
}

func New(testCases []core.TestCase) Model {
	return Model{
		items:  testCases,
		keyMap: DefaultKeyMap(),
		styles: DefaultStyles(true),
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.focused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			m.index = max(m.index-1, 0)
		case key.Matches(msg, m.keyMap.Down):
			m.index = min(m.index+1, len(m.items)-1)
		case key.Matches(msg, m.keyMap.Create):
			core.CreateTestCase(&m.items)
			m.index = len(m.items) - 1
		}
	case tea.MouseReleaseMsg:
		if msg.Button == tea.MouseLeft {
			for i, _ := range m.items {
				if zone.Get("listitem-" + strconv.Itoa(i)).InBounds(msg) {
					m.index = i
					break
				}
			}
		}
	}

	return nil
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(m.styles.Title.Render("main.c"))
	sb.WriteString("\n\n")

	var state StyleState = m.styles.Blurred
	if m.focused {
		state = m.styles.Focused
	}

	for i, testCase := range m.items {
		if i > 0 {
			sb.WriteString("\n\n")
		}

		content := fmt.Sprintf(
			"%s\n%s",
			m.styles.getTitleStyle(testCase.Status).Bold(i == m.index).Render(string(testCase.Status)),
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
