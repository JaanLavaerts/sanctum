package tui

import (
	"errors"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


const (
	Login = "Login"
	Home = "Home"
)

type Model struct {
	currentView string
	masterPassword textinput.Model
	entries []database.Entry
	err       error
}

// TODO: determine initial model based on password recency
func InitialModel() Model {
	m := Model{
		currentView:     Login,
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
	case Login:
		return loginUpdate(msg, m)
	case Home:
		return homeUpdate(msg, m)
	}
	return m, nil
}

func (m Model) View() string {
	switch m.currentView {
	case Login:
		return loginView(m)
	case Home:
		return homeView(m)
	}
	return ""
}

