package tui

import (
	"errors"
	"fmt"

	"github.com/JaanLavaerts/sanctum/crypto"
	"github.com/JaanLavaerts/sanctum/database"
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

				hashed_password, err := database.GetMasterPassword()

				if err != nil {
					m.err = err
				}

				value := crypto.VerifyMasterPassword(m.masterPassword.Value(), hashed_password)

				if err != nil {
					m.err = err
				}

				if value {
					m.currentView = Home
				}
				 	m.err = errors.New("(master password incorrect)")
					m.masterPassword.SetValue("")
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
		"Please provide your master password: %s \n\n%s\n",
		m.err,
		m.masterPassword.View(),
	) + "\n" 
}