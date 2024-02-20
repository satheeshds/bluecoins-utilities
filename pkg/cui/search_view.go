package cui

import (
	"bluecoins-to-splitwise-go/pkg/model"

	"github.com/jroimartin/gocui"
)

type SearchView struct {
	View        *gocui.View
	Transaction *model.BankTransaction
	Name        string
	Text        string
	SearchFn    func(text string) []string
}

const (
	inputViewHeight = 2
)

func (d *SearchView) Create(g *gocui.Gui, x0, y0, x1, y1 int, next func(*gocui.Gui, *gocui.View) error) error {

	input := &InputView{
		Name: d.Name + "_input",
		Text: d.Text,
	}

	if err := input.Create(g, x0, y0, x1, y0+inputViewHeight); err != nil {
		return err
	}
	g.SetCurrentView(input.Name)
	selectablelist := NewSelectableList(d.Name+"_select", x0, y0+inputViewHeight+1, d.SearchFn)
	selectablelist.Layout(g)

	if err := g.SetKeybinding(input.Name, gocui.KeySpace, gocui.ModNone, selectablelist.Update); err != nil {
		return err
	}
	if err := g.SetKeybinding(input.Name, gocui.KeyBackspace, gocui.ModNone, selectablelist.Update); err != nil {
		return err
	}
	if err := g.SetKeybinding(input.Name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(selectablelist.Name)
		return nil
	}); err != nil {
		return err
	}
	g.SetKeybinding(selectablelist.Name, gocui.KeyArrowDown, gocui.ModNone, selectablelist.Down)
	g.SetKeybinding(selectablelist.Name, gocui.KeyArrowUp, gocui.ModNone, selectablelist.Up)
	if err := g.SetKeybinding(selectablelist.Name, gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(input.Name)
		return nil
	}); err != nil {
		return err
	}
	// g.SetKeybinding(m.Name, gocui.KeyEnter, gocui.ModNone, selectablelist.Enter)

	if err := g.SetKeybinding(selectablelist.Name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		selectedText := selectablelist.GetSelected()
		d.Transaction.Description = selectedText
		g.DeleteView(input.Name)
		g.DeleteView(selectablelist.Name)
		return next(g, v)
	}); err != nil {
		return err
	}

	return nil
}
