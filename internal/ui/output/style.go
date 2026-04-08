package output

import "charm.land/lipgloss/v2"

type Styles struct {
	Placeholder lipgloss.Style
	Base        lipgloss.Style
}

func DefaultStyles(isDark bool) Styles {
	return Styles{
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.BrightBlack),
		Base:        lipgloss.NewStyle().Padding(1),
	}
}
