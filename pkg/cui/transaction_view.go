package cui

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"fmt"

	"github.com/jroimartin/gocui"
)

type TransactionView struct {
	Data                       *model.BluecoinsTransactionImport
	Name                       string
	LogHandler                 func(*gocui.View, string)
	Selected                   func(...model.BluecoinsTransactionImport) func(g *gocui.Gui, v *gocui.View) error
	nameView                   *InputView
	categoryView               *SearchView
	termLabel                  *SelectableList
	splitLabel                 *SelectableList
	isTransferView             *BoolView
	transferAccountView        *SearchView
	CategorySearchFn           func(text string) []fmt.Stringer
	AccountSearchfn            func(text string) []fmt.Stringer
	startX, startY, endX, endY int
}

func (t *TransactionView) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	t.startX = x0
	t.startY = y0
	t.endX = x1
	t.endY = y1
	t.nameView = &InputView{
		Name:           t.Name + "Name",
		Text:           t.Data.ItemOrPayee,
		UpdateHandler:  t.UpdateName,
		LogHandler:     t.LogHandler,
		DiscardHandler: t.Discard,
	}
	t.categoryView = &SearchView{
		Name:           t.Name + "Category",
		Text:           t.Data.Category,
		SearchFn:       t.CategorySearchFn,
		UpdateHandler:  t.UpdateCategory,
		LogHandler:     t.LogHandler,
		DiscardHandler: t.Discard,
	}
	terms := []string{"ShortTerm", "MidTerm", "LongTerm"}
	termItems := make([]fmt.Stringer, len(terms))
	for i, term := range terms {
		termItems[i] = model.StringerString(term)
	}
	t.termLabel = &SelectableList{
		Name:              t.Name + "Term",
		Items:             termItems,
		LogHandler:        t.LogHandler,
		SelectedHandler:   t.UpdateLabel,
		InputFocusHandler: t.Discard,
		StartX:            x0,
		StartY:            y0,
	}
	splitLabels := []string{"NotSplit", "SplitEqual", "SplitUnequal"}
	splitItems := make([]fmt.Stringer, len(splitLabels))
	for i, split := range splitLabels {
		splitItems[i] = model.StringerString(split)
	}
	t.splitLabel = &SelectableList{
		Name:              t.Name + "Split",
		Items:             splitItems,
		LogHandler:        t.LogHandler,
		SelectedHandler:   t.UpdateLabel,
		InputFocusHandler: t.Discard,
		StartX:            x0,
		StartY:            y0,
	}
	t.isTransferView = &BoolView{
		Name:       t.Name + "Show",
		LogHandler: t.LogHandler,
		Text:       "Transfer",
		Selected:   t.Show,
	}
	t.transferAccountView = &SearchView{
		Name:           t.Name + "TransferAccount",
		Text:           "",
		SearchFn:       t.AccountSearchfn,
		UpdateHandler:  t.UpdateTransferAccount,
		LogHandler:     t.LogHandler,
		DiscardHandler: t.Discard,
	}

	if err := t.nameView.Layout(g, x0, y0, x1); err != nil {
		return err
	}
	if err := t.FocusView(g, t.nameView.Name); err != nil {
		return err
	}
	return nil
}

func (t *TransactionView) UpdateName(g *gocui.Gui, v *gocui.View) error {
	t.Data.ItemOrPayee, _ = v.Line(0)
	if err := DeleteView(g, t.nameView.Name); err != nil {
		return err
	}

	if err := t.isTransferView.Layout(g, t.startX, t.startY, t.endX, t.startY+2); err != nil {
		return err
	}
	return t.FocusView(g, t.isTransferView.Name)
}

func (t *TransactionView) Show(show bool) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.LogHandler(v, "deleting views")
		DeleteView(g, t.isTransferView.Name)
		t.LogHandler(v, fmt.Sprintf("showing transfer view: %v", show))
		if show {
			transferTxt := "Transfer"
			if t.Data.Type == "e" {
				transferTxt += " to"
			} else {
				transferTxt += " from"
			}
			t.transferAccountView.Name = transferTxt
			return t.transferAccountView.Create(g, t.startX, t.startY, t.endX, t.endY)
		} else {
			return t.categoryView.Create(g, t.startX, t.startY, t.endX, t.endY)
		}
	}
}

func (t *TransactionView) UpdateTransferAccount(account interface{}) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		acc, ok := account.(model.Account)
		if !ok {
			return fmt.Errorf("invalid account type: %T", acc)
		}
		txns := t.Data.GetTransferTransactions(acc)
		return t.Selected(txns...)(g, v)
	}
}

func (t *TransactionView) UpdateCategory(category interface{}) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.LogHandler(v, "updating category")
		category, ok := category.(model.Category)
		if !ok {
			return fmt.Errorf("invalid category type: %T", category)
		}
		t.Data.Category = category.Name
		t.Data.ParentCategory = category.ParentCategory
		t.LogHandler(v, fmt.Sprintf("category: %s, parent: %s", t.Data.Category, t.Data.ParentCategory))
		t.LogHandler(v, fmt.Sprintf("deleting view: %s", t.categoryView.Name))
		if err := t.categoryView.Discard(g, v); err != nil {
			return err
		}
		t.LogHandler(v, fmt.Sprintf("layouting view: %s", t.termLabel.Name))
		if err := t.termLabel.Layout(g); err != nil {
			return err
		}
		t.LogHandler(v, fmt.Sprintf("focusing view: %s", t.termLabel.Name))
		return t.FocusView(g, t.termLabel.Name)
	}
}

func (t *TransactionView) UpdateLabel(lbl interface{}) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		label, ok := lbl.(model.StringerString)
		if !ok {
			return fmt.Errorf("invalid label type: %T", label)
		}
		t.Data.Labels = append(t.Data.Labels, label.String())
		if v.Name() == t.termLabel.Name {
			if err := DeleteView(g, t.termLabel.Name); err != nil {
				return err
			}
			if err := t.splitLabel.Layout(g); err != nil {
				return err
			}
			return t.FocusView(g, t.splitLabel.Name)
		}

		if v.Name() == t.splitLabel.Name {
			if err := DeleteView(g, t.splitLabel.Name); err != nil {
				return err
			}
			return t.Selected(*t.Data)(g, v)
		}
		return nil
	}
}

func (t *TransactionView) Discard(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func (t *TransactionView) FocusView(g *gocui.Gui, viewName string) error {
	if _, err := g.SetCurrentView(viewName); err != nil {
		return err
	}
	return nil
}
