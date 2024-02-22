package cui

import (
	"bluecoins-to-splitwise-go/pkg/bluecoins"
	"bluecoins-to-splitwise-go/pkg/model"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

type MainView struct {
	view                  *gocui.View
	Name                  string
	curTransaction        int
	Transactions          []model.BankTransaction
	blueCoinsTransactions []model.BluecoinsTransactionImport
	Logfile               *os.File
	Verbose               bool
	BluecoinsService      bluecoins.BluecoinsService
}

var (
	sampleItems = []string{
		"Apple",
		"Banana",
		"Cherry",
	}
)

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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

		if err := g.SetKeybinding(m.Name, 'y', gocui.ModNone, m.IncludeTransaction); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(m.Name, 'n', gocui.ModNone, m.Next); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
		// this applies to all views
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
		m.view = v

	}
	v.Clear()
	fmt.Fprintln(v, m.CurrentTransaction().Description)
	fmt.Fprintf(v, "(%d) Add to Bluecoins: (y/n)", m.curTransaction)
	return nil
}

func (m *MainView) Next(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(m.Name); err != nil {
		return err
	}

	if m.curTransaction < len(m.Transactions)-1 {
		m.curTransaction++
		m.AddLog(v, fmt.Sprintf("Next transaction : %d", m.curTransaction))
		m.AddLog(v, fmt.Sprintf("total views : %d", len(g.Views())))
		m.AddLog(v, "----------------------")
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
	for _, transaction := range m.blueCoinsTransactions {
		fmt.Println(transaction.String())
	}
	return nil, nil
}

func (m *MainView) UpdateDescription(transaction interface{}) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		//
		txn, ok := transaction.(model.BluecoinsTransactionImport)
		if !ok {
			return fmt.Errorf("invalid transaction type: %T", transaction)
		}
		m.blueCoinsTransactions = append(m.blueCoinsTransactions, txn)
		return nil
	}
}

func (m *MainView) IncludeTransaction(g *gocui.Gui, v *gocui.View) error {

	descView := &SearchView{
		Text:          m.CurrentTransaction().CleanDescription(),
		Name:          "description",
		UpdateHandler: m.UpdateDescription,
		SearchFn:      m.DescriptionSearch,
		NextHandler:   m.Next,
		LogHandler:    m.AddLog,
	}

	if err := descView.Create(g, 5, 5, 50, 50); err != nil {
		return err
	}

	return nil
}

func (m *MainView) DescriptionSearch(text string) []fmt.Stringer {
	transactions, err := m.BluecoinsService.GetTransactionsImportFormatByDescription(text)
	if err != nil {
		m.AddLog(m.view, fmt.Sprintf("Error getting transactions: %s", err))
		return []fmt.Stringer{}
	}
	return transactions
}

func (m *MainView) AddLog(view *gocui.View, text string) {
	if !m.Verbose {
		return
	}
	viewName := "undefined"
	if view != nil {
		viewName = view.Name()
	}
	fmt.Fprintf(m.Logfile, "[%s] %s\n", viewName, text)
}

func DeleteView(g *gocui.Gui, viewName ...string) error {
	for _, name := range viewName {
		if err := g.DeleteView(name); err != nil {
			return err
		}

		g.DeleteKeybindings(name)
	}
	return nil
}
