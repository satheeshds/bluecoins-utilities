package model

import (
	"fmt"
	"strings"
)

type BluecoinsTransactionImport struct {
	Name           string
	Category       string
	ParentCategory string
	Labels         []string
}

func (t *BluecoinsTransactionImport) ToString() string {
	return fmt.Sprintf("%s|%s|%s|%s", t.Name, t.Category, t.ParentCategory, strings.Join(t.Labels, ","))
}
