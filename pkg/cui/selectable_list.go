package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type SelectableList struct {
	Items             []string
	SelectedIndex     int
	Name              string
	StartX, StartY    int
	SelectedHandler   func(text string) func(g *gocui.Gui, v *gocui.View) error
	InputFocusHandler func(g *gocui.Gui, v *gocui.View) error
	LogHandler        func(*gocui.View, string)
}

func (s *SelectableList) height() int {
	return len(s.Items) + 1
}

func (s *SelectableList) width() int {
	width := 0
	for _, item := range s.Items {
		if len(item) > width {
			width = len(item)
		}
	}
	return width + 1
}

func (s *SelectableList) Layout(g *gocui.Gui) error {
	width := s.width()
	height := s.height()
	v, err := g.SetView(s.Name, s.StartX, s.StartY, s.StartX+width, s.StartY+height)

	if err != nil {
		s.LogHandler(v, fmt.Sprintf("Intialization for: %v", err))
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if err = g.SetKeybinding(s.Name, gocui.KeyArrowDown, gocui.ModNone, s.Down); err != nil {
			return err
		}
		if err = g.SetKeybinding(s.Name, gocui.KeyArrowUp, gocui.ModNone, s.Up); err != nil {
			return err
		}
		if err = g.SetKeybinding(s.Name, gocui.KeyEnter, gocui.ModNone, s.Select); err != nil {
			return err
		}
		if err = g.SetKeybinding(s.Name, gocui.KeyTab, gocui.ModNone, s.InputFocusHandler); err != nil {
			return err
		}
	} else {
		s.LogHandler(v, "Clearing view ")
		v.Clear()
	}
	for _, item := range s.Items {
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
		if err := v.SetCursor(cx, cy+1); err != nil && s.SelectedIndex < len(s.Items)-1 {
			s.SelectedIndex++
			v.MoveCursor(0, 1, false)
		}
	}
	return nil
}

func (s *SelectableList) Select(g *gocui.Gui, v *gocui.View) error {
	s.LogHandler(v, fmt.Sprintf("Selected: %s", s.GetSelected()))
	selectedText := s.GetSelected()
	s.LogHandler(v, fmt.Sprintf("Calling selected handler for %s", selectedText))
	return s.SelectedHandler(selectedText)(g, v)
}

func (s *SelectableList) GetSelected() string {
	if s.SelectedIndex < 0 || s.SelectedIndex >= len(s.Items) {
		return ""
	}

	return s.Items[s.SelectedIndex]
}
