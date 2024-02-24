package cui

import "github.com/jroimartin/gocui"

const (
	inputViewHeight = 2
)

type InputView struct {
	view           *gocui.View
	Name           string
	Text           string
	UpdateHandler  func(g *gocui.Gui, v *gocui.View) error
	DiscardHandler func(g *gocui.Gui, v *gocui.View) error
	LogHandler     func(*gocui.View, string)
}

func (i *InputView) Layout(g *gocui.Gui, x0, y0, width int) error {
	v, err := g.SetView(i.Name, x0, y0, x0+width, y0+inputViewHeight)
	i.LogHandler(v, "Layout")
	if err != nil {
		i.LogHandler(v, "Intialization for: "+err.Error())
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.Title = i.Name
		if i.Text != "" {
			v.Write([]byte(i.Text))
		}
		if err := g.SetKeybinding(i.Name, gocui.KeyTab, gocui.ModNone, i.UpdateHandler); err != nil {
			return err
		}
		if err := g.SetKeybinding(i.Name, gocui.KeyEnter, gocui.ModNone, i.UpdateHandler); err != nil {
			return err
		}
		if err := g.SetKeybinding(i.Name, gocui.KeyCtrlD, gocui.ModNone, i.Discard); err != nil {
			return err
		}

	}
	i.view = v
	return nil
}

func (i *InputView) Discard(g *gocui.Gui, v *gocui.View) error {
	i.LogHandler(v, "Discard")
	return i.DiscardHandler(g, v)
}
