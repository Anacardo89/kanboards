/*
Menu
  |_Project
    |_Label
	|_List
	  |_Card
		|_CheckList
		|_CardLabels
*/

package kanban

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/google/uuid"
)

type Menu struct {
	Projects dll.DLL
}

type Project struct {
	Id     uuid.UUID
	Title  string
	Boards dll.DLL
	Labels dll.DLL
}

type Board struct {
	Id    uuid.UUID
	Title string
	Cards dll.DLL
}

type Label struct {
	Id    uuid.UUID
	Title string
	Color string
}

type Card struct {
	Id          uuid.UUID
	Title       string
	Description string
	CheckList   dll.DLL
	CardLabels  dll.DLL
}

type CheckItem struct {
	Id    uuid.UUID
	Title string
	Check bool
}

// Menu
func StartMenu() *Menu {
	return &Menu{
		Projects: dll.New(),
	}
}

func (m *Menu) AddProject(id uuid.UUID, title string) {
	project := &Project{
		Id:     id,
		Title:  title,
		Boards: dll.New(),
		Labels: dll.New(),
	}
	m.Projects.Append(project)
}

func (m *Menu) RemoveProject(project *Project) error {
	_, err := m.Projects.Remove(project)
	return err
}

// Project
func (p *Project) RenameProject(title string) {
	p.Title = title
}

func (p *Project) AddBoard(id uuid.UUID, title string) {
	board := &Board{
		Id:    id,
		Title: title,
		Cards: dll.New(),
	}
	p.Boards.Append(board)
}

func (p *Project) RemoveBoard(board *Board) error {
	_, err := p.Boards.Remove(board)
	return err
}

func (p *Project) AddLabel(id uuid.UUID, title string, color string) {
	label := &Label{
		Id:    id,
		Title: title,
		Color: color,
	}
	p.Labels.Append(label)
}

func (p *Project) RemoveLabel(label *Label) error {
	_, err := p.Labels.Remove(label)
	return err
}

// Label
func (l *Label) RenameLabel(title string) {
	l.Title = title
}

func (l *Label) ChangeColor(color string) {
	l.Color = color
}

// Board
func (b *Board) RenameBoard(title string) {
	b.Title = title
}

func (b *Board) AddCard(id uuid.UUID, title string, desc string) {
	card := &Card{
		Id:          id,
		Title:       title,
		Description: desc,
		CheckList:   dll.New(),
		CardLabels:  dll.New(),
	}
	b.Cards.Append(card)
}

func (b *Board) RemoveCard(card *Card) error {
	_, err := b.Cards.Remove(card)
	return err
}

// Card
func (c *Card) RenameCard(title string) {
	c.Title = title
}

func (c *Card) AddDescription(description string) {
	c.Description = description
}

func (c *Card) AddCheckItem(id uuid.UUID, title string, done bool) {
	checkItem := &CheckItem{
		Id:    id,
		Title: title,
		Check: done,
	}
	c.CheckList.Append(checkItem)
}

func (c *Card) RemoveCheckItem(checkItem *CheckItem) error {
	_, err := c.CheckList.Remove(checkItem)
	return err
}

func (c *Card) AddLabel(label *Label) {
	c.CardLabels.Append(label)
}

func (c *Card) RemoveLabel(label *Label) error {
	_, err := c.CardLabels.Remove(label)
	return err
}

// CheckItem
func (c *CheckItem) RenameCheckItem(title string) {
	c.Title = title
}

func (c *CheckItem) CheckCheckItem() {
	c.Check = !c.Check
}
