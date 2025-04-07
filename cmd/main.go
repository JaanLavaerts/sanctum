package main

import (
	"log"

	"github.com/JaanLavaerts/sanctum/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}



