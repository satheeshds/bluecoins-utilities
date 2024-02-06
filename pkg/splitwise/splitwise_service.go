package splitwise

import "time"

type SplitwiseService interface {
	GetLastExpenseDate() (time.Time, error)
	SetLastExpenseDate(date time.Time) error
}

type Expense interface {
	Create() error
}
