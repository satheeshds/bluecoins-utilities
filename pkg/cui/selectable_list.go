package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type SelectableList struct {
	Items             []fmt.Stringer
	SelectedIndex     int
	Name              string
	StartX, StartY    int
	SelectedHandler   func(val interface{}) func(g *gocui.Gui, v *gocui.View) error
	InputFocusHandler func(g *gocui.Gui, v *gocui.View) error
	LogHandler        func(*gocui.View, string)
}

func (s *SelectableList) height() int {
	return len(s.Items) + 1
}

func (s *SelectableList) width() int {
	width := 0
	for _, item := range s.Items {
		if len(item.String()) > width {
			width = len(item.String())
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
	s.LogHandler(v, fmt.Sprintf("Up -- SelectedIndex: %d", s.SelectedIndex))
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err == nil {
			s.SelectedIndex--
		} else {
			s.LogHandler(v, fmt.Sprintf("Up else not moving cursor -- err (%v)", err))
			s.InputFocusHandler(g, v)
		}

		s.LogHandler(v, fmt.Sprintf("Up -- Updated index: %d", s.SelectedIndex))
	}
	return nil
}

func (s *SelectableList) Down(g *gocui.Gui, v *gocui.View) error {
	s.LogHandler(v, fmt.Sprintf("Down -- SelectedIndex: %d", s.SelectedIndex))
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err == nil {
			s.SelectedIndex++
		} else {
			s.LogHandler(v, fmt.Sprintf("Down else not moving cursor -- err (%v)", err))
			s.InputFocusHandler(g, v)
		}

		s.LogHandler(v, fmt.Sprintf("Down -- Updated index: %d", s.SelectedIndex))
	}
	return nil
}

func (s *SelectableList) Select(g *gocui.Gui, v *gocui.View) error {
	s.LogHandler(v, fmt.Sprintf("Selected: %s", s.GetSelected()))
	selectedText := s.GetSelected()
	s.LogHandler(v, fmt.Sprintf("Calling selected handler for %s", selectedText))
	return s.SelectedHandler(selectedText)(g, v)
}

func (s *SelectableList) GetSelected() interface{} {
	s.LogHandler(nil, fmt.Sprintf("GetSelected() SelectedIndex: %d", s.SelectedIndex))
	if s.SelectedIndex < 0 || s.SelectedIndex >= len(s.Items) {
		return ""
	}

	return s.Items[s.SelectedIndex]
}
