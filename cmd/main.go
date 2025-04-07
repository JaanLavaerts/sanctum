package main

import (
	"log"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/JaanLavaerts/sanctum/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	database.InitDB()

	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}



