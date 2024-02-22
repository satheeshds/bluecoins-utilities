package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type SearchView struct {
	View          *gocui.View
	Name          string
	Text          string
	SearchFn      func(text string) []fmt.Stringer
	UpdateHandler func(desc interface{}) func(g *gocui.Gui, v *gocui.View) error
	NextHandler   func(g *gocui.Gui, v *gocui.View) error
	inputView     *InputView
	listView      *SelectableList
	LogHandler    func(*gocui.View, string)
}

const (
	inputViewHeight = 2
)

func (s *SearchView) Create(g *gocui.Gui, x0, y0, x1, y1 int) error {
	inputName := s.Name + "_input"
	selectableName := s.Name + "_select"
	s.listView = &SelectableList{
		Name:              selectableName,
		StartX:            x0,
		StartY:            y0 + inputViewHeight + 1,
		SelectedHandler:   s.Selected,
		InputFocusHandler: s.FocusInput,
		Items:             s.SearchFn(s.Text),
		LogHandler:        s.LogHandler,
	}
	s.inputView = &InputView{
		Name:          inputName,
		Text:          s.Text,
		UpdateHandler: s.UpdateList,
		LogHandler:    s.LogHandler,
	}

	if err := s.inputView.Layout(g, x0, y0, x1, y0+inputViewHeight); err != nil {
		return err
	}
	if err := s.listView.Layout(g); err != nil {
		return err
	}
	s.FocusList(g, nil)

	return nil
}

func (s *SearchView) Selected(txn interface{}) func(g *gocui.Gui, v *gocui.View) error {
	transaction, ok := txn.(fmt.Stringer)
	if !ok {
		return func(g *gocui.Gui, v *gocui.View) error {
			return fmt.Errorf("invalid transaction type: %T", txn)
		}
	}
	s.LogHandler(s.View, fmt.Sprintf("selected: %s", transaction))
	return func(g *gocui.Gui, v *gocui.View) error {
		s.LogHandler(v, "calling update handler")
		s.UpdateHandler(txn)(g, v)
		s.LogHandler(v, "deleting views")
		DeleteView(g, s.inputView.Name, s.listView.Name)
		s.LogHandler(v, "calling next handler")
		return s.NextHandler(g, v)
	}
}

func (s *SearchView) FocusInput(g *gocui.Gui, _ *gocui.View) error {
	if _, err := g.SetCurrentView(s.inputView.Name); err != nil {
		return err
	}
	return nil
}

func (s *SearchView) FocusList(g *gocui.Gui, _ *gocui.View) error {
	if _, err := g.SetCurrentView(s.listView.Name); err != nil {
		return err
	}
	return nil
}

func (s *SearchView) UpdateList(g *gocui.Gui, v *gocui.View) error {
	input, err := s.inputView.View.Line(0)
	s.LogHandler(v, fmt.Sprintf("updating list for %s", input))
	if err != nil {
		return err
	}
	s.listView.Items = s.SearchFn(input)
	s.LogHandler(v, fmt.Sprintf("updating list with %v, and layout", s.listView.Items))
	if err = s.listView.Layout(g); err != nil {
		return err
	}
	return s.FocusList(g, v)
}
