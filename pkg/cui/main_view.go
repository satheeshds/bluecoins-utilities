package cui

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

type MainView struct {
	View           *gocui.View
	Name           string
	curTransaction int
	Transactions   []model.BankTransaction
	include        []bool
}

var (
	sampleItems = []string{
		"Apple",
		"Banana",
		"Cherry",
	}
)

func (m *MainView) Create(g *gocui.Gui) error {

	m.include = make([]bool, len(m.Transactions))

	if err := g.SetKeybinding(m.Name, 'y', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		m.include[m.curTransaction] = true
		descView := &SearchView{
			Text:        m.CurrentTransaction().Description,
			Name:        "description",
			Transaction: m.CurrentTransaction(),
			SearchFn: func(text string) []string {
				var matches []string
				for _, items := range sampleItems {
					if strings.Contains(strings.ToLower(items), strings.ToLower(text)) {
						matches = append(matches, items)
					}
				}
				return matches
			},
		}

		if err := descView.Create(g, 5, 5, 50, 50, m.NextTransaction); err != nil {
			return err
		}

		// fmt.Fprintf(v, "Selected: %s", descView.Text)
		// m.Transactions[m.curTransaction].Description = descView.Text
		return nil
	}); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(m.Name, 'n', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		m.include[m.curTransaction] = false
		return m.NextTransaction(g, v)
	}); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(m.Name, gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit

	}); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return nil
}

func (m *MainView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(m.Name, 0, 0, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = m.Name
		if _, err := g.SetCurrentView(m.Name); err != nil {
			return err
		}
	}
	m.View = v
	v.Clear()
	fmt.Fprintln(v, m.CurrentTransaction().Description)
	fmt.Fprintln(v, "Add to Bluecoins: (y/n)")
	return nil
}

func (m *MainView) NextTransaction(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView(m.Name)
	if m.curTransaction < len(m.Transactions)-1 {
		m.curTransaction++
		g.Update(m.Layout)
	} else {
		return gocui.ErrQuit
	}
	return nil
}

func (m *MainView) CurrentTransaction() *model.BankTransaction {
	return &m.Transactions[m.curTransaction]
}

func (m *MainView) GetSelectedTransactions() ([][]string, error) {
	for i, transaction := range m.Transactions {
		if m.include[i] {
			fmt.Println("Included:", transaction.Description)
		} else {
			fmt.Println("Excluded:", transaction.Description)
		}
	}
	return nil, nil
}
