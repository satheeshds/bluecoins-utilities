package cui

import "github.com/jroimartin/gocui"

type InputView struct {
	View *gocui.View
	Name string
	Text string
}

func (i *InputView) Create(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(i.Name, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.Title = i.Name
		if i.Text != "" {
			v.Write([]byte(i.Text))
		}
	}
	i.View = v
	return nil

}
