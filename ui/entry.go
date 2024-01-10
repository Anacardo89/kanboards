package ui

import tea "github.com/charmbracelet/bubbletea"

type modelState int

const (
	menu modelState = iota
	project
	card
)

type Model struct {
	state   modelState
	menu    Menu
	project tea.Model
	card    tea.Model
}

func New() Model {
	return Model{
		state: menu,
		menu:  NewMenu(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case modelState:
		m.state = msg
		switch msg {
		case menu:
			updatedMenu, _ := m.menu.Update(msg)
			m.menu = updatedMenu.(Menu)
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case menu:
		return m.menu.View()
	}
	return ""
}