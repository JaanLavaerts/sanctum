package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)


func homeUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func  homeView(m Model) string {
	return fmt.Sprintf(
		"Your passwords: \n\n%s\n",
		m.entries,
	) + "\n" 
}