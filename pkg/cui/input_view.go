package cui

import "github.com/jroimartin/gocui"

type InputView struct {
	View          *gocui.View
	Name          string
	Text          string
	UpdateHandler func(g *gocui.Gui, v *gocui.View) error
	LogHandler    func(*gocui.View, string)
}

func (i *InputView) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(i.Name, x0, y0, x1, y1)
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
		if err := g.SetKeybinding(i.Name, gocui.KeySpace, gocui.ModNone, i.UpdateHandler); err != nil {
			return err
		}
	}
	i.View = v
	return nil

}
