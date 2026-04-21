package output

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	*viewport.Model

	KeyMap KeyMap
	styles Styles

	focused bool
}

func New() Model {
	model := viewport.New()

	styles := DefaultStyles(true)
	model.SetContent("")

	model.Style = styles.Base

	return Model{
		Model:  &model,
		KeyMap: DefaultKeyMap(),
		styles: styles,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.focused {
		return nil
	}
	model, cmd := m.Model.Update(msg)
	m.Model = &model
	return cmd
}

func (m Model) View() string {
	if m.Model.GetContent() != "" {
		return m.Model.View()
	}

	// render placeholder
	return m.styles.Placeholder.
		Italic(!m.focused).
		Width(m.Width()).Height(m.Height()).
		Align(lipgloss.Center, lipgloss.Center).
		Render("Run (^r) the testcase to see the output")
}

func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) SetStyles(styles Styles) {
	m.styles = styles
	m.Style = styles.Base
}
