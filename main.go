package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
	"github.com/thecomputerm/lazycph/internal/workspace"
)

func main() {
	zone.NewGlobal()
	p := tea.NewProgram(workspace.New())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
