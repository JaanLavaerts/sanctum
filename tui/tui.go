package tui

import (
	"errors"
	"log"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


const (
	Welcome = "Welcome"
	Login = "Login"
	Home = "Home"
)

type Model struct {
	currentView string
	masterPassword textinput.Model
	entries []database.Entry
	err       error
}

func InitialModel() Model {
	currentView := Login
	hashed_password, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}
	if len(hashed_password) == 0 {
		currentView = Welcome
	}
	
	m := Model{
		currentView:     currentView,
		masterPassword: textinput.New(),
		err:           errors.New(""),
	}
	m.masterPassword.Focus()
	m.masterPassword.EchoMode = textinput.EchoPassword
	m.masterPassword.EchoCharacter = ('•')
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case Welcome:
		return welcomeUpdate(msg, m)
	case Login:
		return loginUpdate(msg, m)
	case Home:
		return homeUpdate(msg, m)
	}
	return m, nil
}

func (m Model) View() string {
	switch m.currentView {
	case Welcome:
		return welcomeView(m)
	case Login:
		return loginView(m)
	case Home:
		return homeView(m)
	}
	return ""
}

