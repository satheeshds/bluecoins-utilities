package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type SelectableList struct {
	items          []string
	SelectedIndex  int
	Name           string
	StartX, StartY int
	SearchFn       func(text string) []string
}

func NewSelectableList(name string, startX, startY int, searchFn func(text string) []string) *SelectableList {
	return &SelectableList{
		Name:     name,
		StartX:   startX,
		StartY:   startY,
		SearchFn: searchFn,
		items:    searchFn(""),
	}
}

func (s *SelectableList) height() int {
	return len(s.items) + 1
}

func (s *SelectableList) width() int {
	width := 0
	for _, item := range s.items {
		if len(item) > width {
			width = len(item)
		}
	}
	return width + 1
}

func (s *SelectableList) Layout(g *gocui.Gui) error {
	// maxX, maxY := g.Size()
	width := s.width()
	height := s.height()
	v, err := g.SetView(s.Name, s.StartX, s.StartY, s.StartX+width, s.StartY+height)
	// v, err := g.SetView(s.Name, maxX/2-7, maxY/2-2, maxX/2+7, maxY/2+2)
	// v.Clear()

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

	}
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	v.Clear()
	for _, item := range s.items {
		fmt.Fprintln(v, item)
	}
	return nil
}

func (s *SelectableList) Up(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			s.SelectedIndex--
			v.MoveCursor(0, -1, false)
		}
	}
	return nil
}

func (s *SelectableList) Down(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil && s.SelectedIndex < len(s.items)-1 {
			s.SelectedIndex++
			v.MoveCursor(0, 1, false)
		}
	}
	return nil
}

func (s *SelectableList) Select(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func (s *SelectableList) Update(g *gocui.Gui, v *gocui.View) error {
	input := v.BufferLines()[0]
	s.items = s.SearchFn(input)
	s.Layout(g)
	g.SetCurrentView(s.Name)
	return nil
}

func (s *SelectableList) GetSelected() string {
	if s.SelectedIndex < 0 || s.SelectedIndex >= len(s.items) {
		return ""
	}

	return s.items[s.SelectedIndex]
}
