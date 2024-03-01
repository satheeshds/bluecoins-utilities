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
	m.AddLog(v, "Layout")
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
	fmt.Fprintf(v, "%-20s:%s\n", "Type", curr.TransactionType.String())
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
		m.PopulateTransactionFields(&txn)
		if txn.IsTransfer() {
			transferAccountView := &SearchView{
				Text:     "",
				Name:     "transferAccount",
				SearchFn: m.AccountSearch,
				UpdateHandler: func(account interface{}) func(g *gocui.Gui, v *gocui.View) error {
					return func(g *gocui.Gui, v *gocui.View) error {
						acc, ok := account.(model.Account)
						if !ok {
							return fmt.Errorf("invalid account type: %T", acc)
						}

						txns := txn.GetTransferTransactions(acc)
						return m.AddTransaction(txns...)(g, v)
					}

				},
				LogHandler: m.AddLog,
				DiscardHandler: func(g *gocui.Gui, v *gocui.View) error {
					return nil
				},
			}
			return transferAccountView.Create(g, 5, 10, 50, 50)
		} else {
			return m.AddTransaction(txn)(g, v)
		}
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

	return descView.Create(g, 5, 10, 50, 50)
}

func (m *MainView) PopulateTransactionFields(txn *model.BluecoinsTransactionImport) error {
	current := m.CurrentTransaction()
	txn.AccountType = "Bank"
	txn.Account = "SBI Technopark"
	txn.Date = current.Date.Format("01/02/2006")
	txn.Amount = fmt.Sprintf("%f", current.Amount)
	if current.TransactionType == model.Credit {
		txn.Type = "i"
	} else {
		txn.Type = "e"
	}
	return nil
}

func (m *MainView) PopulateTransaction(g *gocui.Gui, v *gocui.View) error {
	m.AddLog(v, "PopulateTransaction")
	txnImport := &model.BluecoinsTransactionImport{}
	m.PopulateTransactionFields(txnImport)

	transView := &TransactionView{
		Data:             txnImport,
		Name:             "transaction",
		LogHandler:       m.AddLog,
		Selected:         m.AddTransaction,
		CategorySearchFn: m.CategorySearch,
		AccountSearchfn:  m.AccountSearch,
	}

	return transView.Layout(g, 5, 10, 50, 50)
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
