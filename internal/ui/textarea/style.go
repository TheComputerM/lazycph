package textarea

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"
)

func DefaultStyles(isDark bool) textarea.Styles {
	styles := textarea.DefaultStyles(isDark)
	accent := lipgloss.Color("63")

	styles.Focused.Prompt = styles.Focused.Prompt.Foreground(accent)

	styles.Blurred.LineNumber = styles.Focused.CursorLineNumber
	styles.Blurred.CursorLineNumber = styles.Focused.CursorLineNumber.Italic(true)
	styles.Focused.LineNumber = styles.Focused.CursorLineNumber

	styles.Focused.CursorLineNumber = lipgloss.NewStyle().Foreground(accent)

	return styles
}
