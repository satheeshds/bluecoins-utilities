package cui

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"fmt"

	"github.com/jroimartin/gocui"
)

type TransferView struct {
	Name               string
	Transaction        model.BluecoinsTransactionImport
	LogHandler         func(*gocui.View, string)
	Selected           func(*model.BluecoinsTransactionImport) func(g *gocui.Gui, v *gocui.View) error
	showView           *BoolView
	accountView        *SearchView
	counterTransaction *model.BluecoinsTransactionImport
	AccountSearchfn    func(text string) []fmt.Stringer
}

func (t *TransferView) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	transferText := "Transfer"
	if t.Transaction.Type == "e" {
		transferText += " to"
	} else {
		transferText += " from"
	}
	t.accountView = &SearchView{
		Name:          t.Name + "Account",
		Text:          transferText,
		SearchFn:      t.AccountSearchfn,
		UpdateHandler: t.SelectedAccount,
		LogHandler:    t.LogHandler,
		DiscardHandler: func(g *gocui.Gui, v *gocui.View) error {
			return nil
		},
	}
	t.showView = &BoolView{
		Name:       t.Name + "Show",
		LogHandler: t.LogHandler,
		Text:       "Transfer",
		Selected:   t.Show,
	}
	if err := t.showView.Layout(g, x0, y0, x1, y1); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(t.showView.Name); err != nil {
		return err
	}

	return nil
}

func (t *TransferView) Show(show bool) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.LogHandler(v, "deleting views")
		DeleteView(g, t.showView.Name)
		t.LogHandler(v, fmt.Sprintf("showing transfer view: %v", show))
		if show {
			if err := t.accountView.Create(g, 5, 5, 50, 50); err != nil {
				return err
			}
			t.LogHandler(v, "setting current view to account view")
			if _, err := g.SetCurrentView(t.accountView.inputView.Name); err != nil {
				return err
			}
		} else {
			return t.Selected(t.counterTransaction)(g, v)
		}
		return nil
	}
}

func (t *TransferView) SelectedAccount(account interface{}) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		DeleteView(g, t.showView.Name)
		acc, ok := account.(model.Account)
		if !ok {
			return fmt.Errorf("invalid account type: %T", acc)
		}
		t.counterTransaction = &model.BluecoinsTransactionImport{
			Type:           "t",
			Date:           t.Transaction.Date,
			Amount:         t.Transaction.Amount,
			ItemOrPayee:    t.Transaction.ItemOrPayee,
			Category:       "(Transfer)",
			ParentCategory: "(Transfer)",
			Account:        acc.Name,
			AccountType:    acc.TypeName,
		}
		return t.Selected(t.counterTransaction)(g, v)
	}
}
