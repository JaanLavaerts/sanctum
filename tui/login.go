package tui

import (
	"fmt"
	"os"

	"github.com/JaanLavaerts/sanctum/crypto"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)


func loginUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
			case tea.KeyEnter:

				hash, err := os.ReadFile("data/db.txt")

				if err != nil {
					m.err = err
				}

				value, err := crypto.VerifyMasterPassword(string(hash), m.masterPassword.Value())

				if err != nil {
					m.err = err
				}

				if value {
					m.currentView = Home
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

func  loginView(m Model) string {
	return fmt.Sprintf(
		"Please provide your master password:\n\n%s\n\n%s",
		m.masterPassword.View(),
		"(esc to quit)",
	) + "\n"
}