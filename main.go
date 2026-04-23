package main

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/thecomputerm/lazycph/internal/app"
)

func main() {
	currentDirectory, _ := os.Getwd()
	p := tea.NewProgram(app.New(currentDirectory))
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
