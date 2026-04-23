package workspace

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/core"
)

func (m Model) selectFileCmd() tea.Msg {
	return core.NavigateMsg{Path: filepath.Dir(m.filePath)}
}
