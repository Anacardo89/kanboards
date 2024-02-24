package ui

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cursorPos int

const (
	titlePos cursorPos = iota
	descPos
	checkPos
	labelPos
)

// Implements tea.Model
type Card struct {
	card      *kanban.Card
	checklist list.Model
	labels    list.Model
	cursor    cursorPos
	textinput textinput.Model
	textarea  textarea.Model
	flag      actionFlag
}

func (c Card) Init() tea.Cmd {
	return nil
}

func (c Card) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		c.setLists()
		return c, nil
	case tea.KeyMsg:
		if c.textinput.Focused() {
			cmd = c.inputFocused(msg)
			return c, cmd
		}
		if c.textarea.Focused() {
			switch msg.String() {
			case "esc":
				c.textarea.Blur()
			}
			c.textarea, cmd = c.textarea.Update(msg)
			return c, cmd
		}
		cmd = c.keyPress(msg)
		return c, cmd
	}
	if c.cursor == checkPos {
		c.checklist, cmd = c.checklist.Update(msg)
	} else if c.cursor == labelPos {
		c.labels, cmd = c.labels.Update(msg)
	}
	return c, cmd
}

func (c Card) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	return c.cardView()

}

// called by model_selector
func OpenCard(kc *kanban.Card) Card {
	c := Card{
		card:      kc,
		textinput: textinput.New(),
		textarea:  textarea.New(),
		cursor:    0,
	}
	c.setInput()
	c.setTxtArea()
	c.setLists()
	return c
}

func (c *Card) UpdateCard() {
	c.setLists()
	c.setTxtArea()
}

func (c *Card) getCheckItem() *kanban.CheckItem {
	if c.card.CheckList.Length() == 0 {
		return nil
	}
	ci, err := c.card.CheckList.GetAt(c.checklist.Cursor())
	if err != nil {
		log.Println(err)
		return nil
	}
	return ci.(*kanban.CheckItem)
}

func (c *Card) getCardLabel() *kanban.Label {
	if c.card.CardLabels.Length() == 0 {
		return nil
	}
	l, err := c.card.CardLabels.GetAt(c.labels.Cursor())
	if err != nil {
		log.Println(err)
		return nil
	}
	return l.(*kanban.Label)
}

// Update
func (c *Card) inputFocused(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		c.textinput.SetValue("")
		c.textinput.Blur()
		c.flag = none
		return nil
	case "enter":
		c.txtInputEnter()
	}
	c.textinput, cmd = c.textinput.Update(msg)
	return cmd
}

func (c *Card) txtInputEnter() {
	if c.textinput.Value() == "" {
		return
	}
	switch c.cursor {
	case 0:
		c.card.RenameCard(c.textinput.Value())
	case 2:
		c.card.AddCheckItem(c.textinput.Value())
	}
	c.setLists()
	c.textinput.SetValue("")
	c.textinput.Blur()
	c.flag = none
}

func (c *Card) keyPress(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "left":
		c.handleMoveLeft()
	case "right":
		c.handleMoveRight()
	case "esc":
		return func() tea.Msg { return upProject }
	case "enter":
		c.setInput()
		switch c.cursor {
		case titlePos:
			c.textinput.Placeholder = "Card Title"
			c.textinput.Focus()
		case descPos:
			c.textarea.Focus()
		case checkPos:
			checkitem := c.getCheckItem()
			checkitem.CheckCheckItem()
			c.setLists()
		case labelPos:
			return func() tea.Msg { return labelState }
		}
	case "n":
		if c.cursor == checkPos {
			c.textinput.Placeholder = "CheckItem Title"
			c.textinput.Focus()
		}
		return nil
	case "d":
		if c.cursor == checkPos || c.cursor == labelPos {
			c.handleDelete()
		}
	}
	return nil
}

// actions
// movement
func (c *Card) handleMoveRight() {
	if c.cursor == labelPos {
		c.cursor = titlePos
	} else {
		c.cursor++
	}
}

func (c *Card) handleMoveLeft() {
	if c.cursor == titlePos {
		c.cursor = labelPos
	} else {
		c.cursor--
	}
}

// delete
func (c *Card) handleDelete() {
	switch c.cursor {
	case checkPos:
		ci, err := c.card.CheckList.GetAt(c.checklist.Cursor())
		if err != nil {
			log.Println(err)
		}
		c.card.RemoveCheckItem(ci.(*kanban.CheckItem))
	case labelPos:
		l, err := c.card.CardLabels.GetAt(c.labels.Cursor())
		if err != nil {
			log.Println(err)
		}
		c.card.RemoveLabel(l.(*kanban.Label))
	}
	c.setLists()
}

// View
func (c *Card) cardView() string {
	var inputStyled = ""
	if c.textinput.Focused() {
		inputStyled = InputFieldStyle.Render(c.textinput.View())
	}
	cardStyled := c.renderCard()
	return lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			cardStyled,
			inputStyled,
		),
	)
}

func (c *Card) renderCard() string {
	emptyLine := ""
	titleStyled := "Title"
	descriptionStyled := "Description"
	txtareaStyled := TextAreaStyle.Render(c.textarea.View())
	checklistStyled := ListStyle.Render(c.checklist.View())
	cardlabelsStyled := ListStyle.Render(c.labels.View())
	switch c.cursor {
	case titlePos:
		titleStyled = SelectedTxtStyle.Render(titleStyled)
	case descPos:
		descriptionStyled = SelectedTxtStyle.Render(descriptionStyled)
	case checkPos:
		checklistStyled = SelectedListStyle.Render(c.checklist.View())
	case labelPos:
		cardlabelsStyled = SelectedListStyle.Render(c.labels.View())
	}
	listsStyled := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checklistStyled,
		cardlabelsStyled,
	)
	return CardStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyled,
		c.card.Title,
		emptyLine,
		descriptionStyled,
		txtareaStyled,
		listsStyled,
	))
}

// bubbles
// textinput
func (c *Card) setInput() {
	c.textinput.Prompt = ": "
	c.textinput.CharLimit = 120
	c.textinput.Cursor.Blink = true
}

// textarea
func (c *Card) setTxtArea() {
	c.textarea.Prompt = ""
	c.textarea.Placeholder = "Card Description"
	c.textarea.ShowLineNumbers = true
	c.textarea.Cursor.Blink = true
	c.textarea.SetValue(c.card.Description)
}

// list
var checklistDelegate = NewCheckListDelegate()

func (c *Card) setLists() {
	c.setCheckList()
	c.setCardLabels()
}

func (c *Card) setCheckList() {
	var checklistItems []list.Item
	cl := list.New([]list.Item{}, checklistDelegate, ws.width/2, ws.height/3+1)
	cl.SetShowHelp(false)
	cl.Title = "Checklist"
	cl.InfiniteScrolling = true
	for i := 0; i < c.card.CheckList.Length(); i++ {
		ci, _ := c.card.CheckList.GetAt(i)
		item := Item{
			title: ci.(*kanban.CheckItem).Title,
		}
		if ci.(*kanban.CheckItem).Check {
			item.description = "1"
		} else {
			item.description = "0"
		}
		checklistItems = append(checklistItems, item)
	}
	cl.SetItems(checklistItems)
	c.checklist = cl
}

func (c *Card) setCardLabels() {
	var labelItems []list.Item
	ll := list.New([]list.Item{}, NewLabelListDelegate(), ws.width/2, ws.height/3+1)
	ll.SetShowHelp(false)
	ll.Title = "Card Labels"
	ll.InfiniteScrolling = true
	for i := 0; i < c.card.CardLabels.Length(); i++ {
		l, _ := c.card.CardLabels.GetAt(i)
		item := Item{
			title: l.(*kanban.Label).Title,
		}
		labelItems = append(labelItems, item)
	}
	ll.SetItems(labelItems)
	c.labels = ll
}