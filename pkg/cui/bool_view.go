package cui

import "github.com/jroimartin/gocui"

type BoolView struct {
	Name       string
	LogHandler func(*gocui.View, string)
	Text       string
	Selected   func(bool) func(g *gocui.Gui, v *gocui.View) error
}

func (b *BoolView) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(b.Name, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = b.Name
		if _, err := g.SetCurrentView(b.Name); err != nil {
			return err
		}
		if err := g.SetKeybinding(b.Name, 'y', gocui.ModNone, b.Selected(true)); err != nil && err != gocui.ErrQuit {
			return err
		}
		if err := g.SetKeybinding(b.Name, 'n', gocui.ModNone, b.Selected(false)); err != nil && err != gocui.ErrQuit {
			return err
		}
	}
	v.Clear()
	v.Write([]byte(b.Text + " (y/n)"))
	return nil
}
