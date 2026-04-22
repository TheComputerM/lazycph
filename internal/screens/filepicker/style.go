package filepicker

import "charm.land/lipgloss/v2"

type Styles struct {
	Title lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Bold(true).MarginBottom(1),
	}
}
