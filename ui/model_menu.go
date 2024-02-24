package ui

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Implements tea.Model
type Menu struct {
	menu      *kanban.Menu
	cursor    int
	list      list.Model
	textinput textinput.Model
	empty     bool
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		m.setList()
		return m, nil
	case tea.KeyMsg:
		if m.textinput.Focused() {
			cmd = m.inputFocused(msg)
			return m, cmd
		}
		cmd = m.keyPress(msg)
		return m, cmd
	}
	m.list, cmd = m.list.Update(msg)
	m.cursor = m.list.Cursor()
	return m, cmd
}

func (m Menu) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	if m.empty {
		return m.viewEmpty()
	}
	return m.viewMenu()
}

// **************************************
func TestData() Menu {
	return Menu{
		cursor:    0,
		menu:      kanban.TestData(),
		textinput: textinput.New(),
	}
}

// **************************************

// called by model_selector
func NewMenu() Menu {
	m := Menu{
		menu:      kanban.StartMenu(),
		textinput: textinput.New(),
		empty:     true,
	}
	m.setTxtInput()
	setMenuDelegate()
	m.setList()
	return m
}

func (m *Menu) UpdateMenu() {
	m.setList()
}

func (m *Menu) getProject() *kanban.Project {
	if m.empty {
		return nil
	}
	project, err := m.menu.Projects.GetAt(m.cursor)
	if err != nil {
		log.Println(err)
	}
	return project.(*kanban.Project)
}

// Update
func (m *Menu) keyPress(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	if m.textinput.Focused() {
		m.textinput, cmd = m.textinput.Update(msg)
		return cmd
	}
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "enter":
		if m.empty {
			return nil
		}
		return func() tea.Msg { return projectState }
	case "n":
		return m.textinput.Focus()
	case "d":
		m.deleteProject()
		return nil
	}
	return nil
}

func (m *Menu) inputFocused(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		m.textinput.SetValue("")
		m.textinput.Blur()
	case "enter":
		m.txtInputEnter()
	}
	return nil
}

func (m *Menu) txtInputEnter() {
	if m.textinput.Value() == "" {
		return
	}
	m.menu.AddProject(m.textinput.Value())
	m.empty = false
	m.setList()
	m.textinput.SetValue("")
	m.textinput.Blur()
	m.cursor = m.list.Cursor()
}

// actions
func (m *Menu) deleteProject() {
	var err error
	if m.empty {
		return
	}
	project, err := m.menu.Projects.GetAt(m.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	err = m.menu.RemoveProject(project.(*kanban.Project))
	if err != nil {
		log.Println(err)
	}
	if m.menu.Projects.Length() == 0 {
		m.empty = true
	}
	m.setList()
	m.cursor = m.list.Cursor()
}

// View
func (m *Menu) viewEmpty() string {
	var (
		bottomLines string
		inputStyled string
	)
	emptyTxtStyled := EmptyStyle.Render(
		"No projects.\n\nPress 'n' to create a new Project Board\nor 'q' to quit",
	)
	if m.textinput.Focused() {
		_, h := lipgloss.Size(emptyTxtStyled)
		for i := 0; i < ws.height-h-h/2; i++ {
			bottomLines += "\n"
		}
		inputStyled = InputFieldStyle.Render(m.textinput.View())
	}
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Center, lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Center,
			emptyTxtStyled,
			bottomLines,
			inputStyled,
		),
	)
}

func (m *Menu) viewMenu() string {
	var (
		bottomLines string
		inputStyled string
	)
	menuStyled := ListStyle.Render(m.list.View())
	if m.textinput.Focused() {
		inputStyled = InputFieldStyle.Render(m.textinput.View())
	}
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Left, lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left,
			menuStyled,
			bottomLines,
			inputStyled,
		),
	)
}

// bubbles
// text input
func (m *Menu) setTxtInput() {
	m.textinput.Prompt = ": "
	m.textinput.CharLimit = 120
	m.textinput.Cursor.Blink = true
	m.textinput.Placeholder = "Project Title"
}

// list
var menuDelegate = list.NewDefaultDelegate()

func setMenuDelegate() {
	menuDelegate.ShowDescription = false
	menuDelegate.SetSpacing(0)
	menuDelegate.Styles.NormalTitle.Foreground(WHITE)
	menuDelegate.Styles.SelectedTitle.Foreground(YELLOW).
		Border(lipgloss.HiddenBorder(), false, false, false, true)
}

func (m *Menu) setList() {
	var menuItems []list.Item
	l := list.New([]list.Item{}, menuDelegate, ws.width/3, ws.height-6)
	l.SetShowHelp(false)
	l.Title = "Projects"
	l.InfiniteScrolling = true
	for i := 0; i < m.menu.Projects.Length(); i++ {
		project, _ := m.menu.Projects.GetAt(i)
		item := Item{
			title: project.(*kanban.Project).Title,
		}
		menuItems = append(menuItems, item)
	}
	l.SetItems(menuItems)
	m.list = l
}