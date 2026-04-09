package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/app"
)

func main() {
	p := tea.NewProgram(app.New())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
