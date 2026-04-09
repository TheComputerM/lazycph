package list

import (
	"charm.land/lipgloss/v2"
	"github.com/thecomputerm/lazycph/internal/core"
)

type StyleState struct {
	SelectedItem lipgloss.Style
	ItemDesc     lipgloss.Style
}

type Styles struct {
	List  lipgloss.Style
	Title lipgloss.Style
	Item  lipgloss.Style

	Focused StyleState
	Blurred StyleState

	PendingItemTitle lipgloss.Style
	CorrectItemTitle lipgloss.Style
	WrongItemTitle   lipgloss.Style
	ErrorItemTitle   lipgloss.Style
}

func DefaultStyles(isDark bool) Styles {
	width := 25
	accent := lipgloss.Color("63")

	selectedItemStyle := lipgloss.NewStyle().
		PaddingLeft(1).
		Width(width).
		Border(lipgloss.ThickBorder(), false, false, false, true)

	return Styles{
		List:  lipgloss.NewStyle().Width(width).Margin(0, 1),
		Title: lipgloss.NewStyle().MaxWidth(width).Padding(0, 1).Background(accent),

		Item: lipgloss.NewStyle().PaddingLeft(2).Width(width),

		Focused: StyleState{
			SelectedItem: selectedItemStyle.BorderForeground(accent),
			ItemDesc:     lipgloss.NewStyle().Foreground(lipgloss.BrightWhite),
		},

		Blurred: StyleState{
			SelectedItem: selectedItemStyle.BorderForeground(lipgloss.White),
			ItemDesc:     lipgloss.NewStyle(),
		},

		PendingItemTitle: lipgloss.NewStyle(),
		CorrectItemTitle: lipgloss.NewStyle().Foreground(lipgloss.BrightGreen),
		WrongItemTitle:   lipgloss.NewStyle().Foreground(lipgloss.BrightRed),
		ErrorItemTitle:   lipgloss.NewStyle().Foreground(lipgloss.BrightYellow),
	}
}

func (m *Model) GetWidth() int {
	return 25 + m.styles.List.GetHorizontalFrameSize()
}

func (s *Styles) itemTitleStyle(status core.TestCaseStatus) lipgloss.Style {
	switch status {
	case core.TestCaseStatusCorrect:
		return s.CorrectItemTitle
	case core.TestCaseStatusWrong:
		return s.WrongItemTitle
	case core.TestCaseStatusError:
		return s.ErrorItemTitle
	default:
		return s.PendingItemTitle
	}
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) SetStyles(styles Styles) {
	m.styles = styles
}
