package splitwise

import "time"

type SplitwiseService interface {
	GetLastExpenseDate(accountId int) (time.Time, error)
	SetLastExpenseDate(accountId int, date time.Time)
	Close()
}

type Expense interface {
	Create() error
}
