package tui

import (
	"fmt"

	"github.com/JaanLavaerts/sanctum/database"
	tea "github.com/charmbracelet/bubbletea"
)


func welcomeUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
			case tea.KeyEnter:

				value, err := database.InserMasterPassword(m.masterPassword.Value())

				if err != nil {
					m.err = err
				}

				if value == 1 {
					m.currentView = Login
					m.masterPassword.SetValue("")
				}
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.masterPassword, cmd = m.masterPassword.Update(msg)

	return m, cmd
}

func  welcomeView(m Model) string {
	return fmt.Sprintf(
		"Welcome! Please provide a master password: %s \n\n%s\n",
		m.err,
		m.masterPassword.View(),
	) + "\n" 
}