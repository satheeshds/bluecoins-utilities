package cui

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"fmt"

	"github.com/jroimartin/gocui"
)

type TransactionView struct {
	Data                       model.BluecoinsTransactionImport
	Name                       string
	LogHandler                 func(*gocui.View, string)
	Selected                   func(...model.BluecoinsTransactionImport) func(g *gocui.Gui, v *gocui.View) error
	nameView                   *InputView
	categoryView               *SearchView
	termLabel                  *SelectableList
	splitLabel                 *SelectableList
	transferView               *TransferView
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

	t.transferView = &TransferView{
		Name:            t.Name + "Transfer",
		LogHandler:      t.LogHandler,
		AccountSearchfn: t.AccountSearchfn, //TODO: change this to AccountSearchFn
		Transaction:     t.Data,
		Selected:        t.UpdateTransfer,
	}
	if err := t.nameView.Layout(g, x0, y0, x1); err != nil {
		return err
	}
	if err := t.FocusView(g, t.nameView.Name); err != nil {
		return err
	}
	return nil
}

// func (t *TransactionView) ValidateAndSubmit(g *gocui.Gui, v *gocui.View) error {
// 	if t.Data.Name != "" {
// 		return t.Selected(t.Data)(g, v)
// 	}
// 	return nil
// }

func (t *TransactionView) UpdateName(g *gocui.Gui, v *gocui.View) error {
	t.Data.ItemOrPayee, _ = v.Line(0)
	if err := DeleteView(g, t.nameView.Name); err != nil {
		return err
	}
	if err := t.transferView.Layout(g, t.startX, t.startY, t.endX, t.endY); err != nil {
		return err
	}
	return t.FocusView(g, t.transferView.showView.Name)
}

func (t *TransactionView) UpdateTransfer(counter *model.BluecoinsTransactionImport) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if counter == nil {
			if err := t.categoryView.Create(g, t.startX, t.startY, t.endX, t.endY); err != nil {
				return err
			}
			return t.FocusView(g, t.categoryView.inputView.Name)
		} else {
			t.Data.Type = "t"
			t.Data.Category = "(Transfer)"
			t.Data.ParentCategory = "(Transfer)"
			return t.Selected(t.Data, *counter)(g, v)
		}

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
			return t.Selected(t.Data)(g, v)
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
