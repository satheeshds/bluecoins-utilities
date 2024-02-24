package model

import "fmt"

type Category struct {
	Name           string
	ParentCategory string
}

func (c Category) String() string {
	return fmt.Sprintf("%s > %s", c.Name, c.ParentCategory)
}
