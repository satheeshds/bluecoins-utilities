package model

import "fmt"

type Account struct {
	ID       int
	Name     string
	TypeName string
}

func (a Account) String() string {
	return fmt.Sprintf("%s > %s", a.TypeName, a.Name)
}
