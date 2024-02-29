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
	curr := m.CurrentTransaction()
	fmt.Fprintf(v, "%-20s:(%d/%d)\n", "Transaction", m.curTransaction+1, len(m.Transactions))
	fmt.Fprintf(v, "%-20s:%s\n", "Description", curr.Description)
	fmt.Fprintf(v, "%-20s:%v\n", "Date", curr.Date)
	fmt.Fprintf(v, "%-20s:%f\n", "Amount", curr.Amount)
	fmt.Fprintf(v, "%-20s: (y/n)", "Add to Bluecoins")
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

func (m *MainView) GetSelectedTransactions() []model.BluecoinsTransactionImport {
	return m.blueCoinsTransactions
}

func (m *MainView) UpdateTransaction(transaction interface{}) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		txn, ok := transaction.(model.BluecoinsTransactionImport)
		if !ok {
			return fmt.Errorf("invalid transaction type: %T", transaction)
		}
		return m.AddTransaction(txn)(g, v)
	}
}

func (m *MainView) IncludeTransaction(g *gocui.Gui, v *gocui.View) error {

	descView := &SearchView{
		Text:           m.CurrentTransaction().CleanDescription(),
		Name:           "description",
		UpdateHandler:  m.UpdateTransaction,
		SearchFn:       m.DescriptionSearch,
		LogHandler:     m.AddLog,
		DiscardHandler: m.PopulateTransaction,
	}

	if err := descView.Create(g, 5, 5, 50, 50); err != nil {
		return err
	}

	return nil
}

func (m *MainView) PopulateTransaction(g *gocui.Gui, v *gocui.View) error {
	m.AddLog(v, "PopulateTransaction")
	current := m.CurrentTransaction()

	txnImport := &model.BluecoinsTransactionImport{
		Date:        current.Date.Format("01/02/2006"),
		AccountType: "Bank",
		Account:     "SBI Technopark",
		Amount:      fmt.Sprintf("%f", current.Amount),
	}
	if current.TransactionType == model.Credit {
		txnImport.Type = "i"
	} else {
		txnImport.Type = "e"
	}

	transView := &TransactionView{
		Data:             *txnImport,
		Name:             "transaction",
		LogHandler:       m.AddLog,
		Selected:         m.AddTransaction,
		CategorySearchFn: m.CategorySearch,
		AccountSearchfn:  m.AccountSearch,
	}

	if err := transView.Layout(g, 5, 5, 50, 50); err != nil {
		return err
	}
	return nil
}

func (m *MainView) DescriptionSearch(text string) []fmt.Stringer {
	var stringer []fmt.Stringer
	transactions, err := m.BluecoinsService.GetTransactionsImportFormatByDescription(text)
	if err != nil {
		m.AddLog(m.view, fmt.Sprintf("Error getting transactions: %s", err))
		return stringer
	}
	for _, txn := range transactions {
		stringer = append(stringer, txn)
	}
	return stringer
}

func (m *MainView) CategorySearch(text string) []fmt.Stringer {
	var stringer []fmt.Stringer
	categories, err := m.BluecoinsService.GetCategories(text)
	if err != nil {
		m.AddLog(m.view, fmt.Sprintf("Error getting categories: %s", err))
		return stringer
	}
	for _, cat := range categories {
		stringer = append(stringer, cat)
	}
	return stringer
}

func (m *MainView) AccountSearch(text string) []fmt.Stringer {
	var stringer []fmt.Stringer
	accounts, err := m.BluecoinsService.GetAccountsBySearch(text)
	if err != nil {
		m.AddLog(m.view, fmt.Sprintf("Error getting accounts: %s", err))
		return stringer
	}
	for _, acc := range accounts {
		stringer = append(stringer, acc)
	}
	return stringer
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

func (m *MainView) AddTransaction(selected ...model.BluecoinsTransactionImport) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		m.blueCoinsTransactions = append(m.blueCoinsTransactions, selected...)
		return m.Next(g, v)
	}
}
