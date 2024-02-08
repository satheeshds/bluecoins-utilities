package model

import "time"

type SplitStatus int

const (
	NotSplit SplitStatus = iota
	SplitEqual
	SplitUnequal
	Undefined
)

// Transaction represents a Bluecoins transaction
type Transaction struct {
	ID          string
	Date        time.Time
	Amount      float32
	Category    int
	Description string
	Labels      []string
}

func (t *Transaction) GetSplitStatus() SplitStatus {
	for _, label := range t.Labels {
		if label == "SplitEqual" {
			return SplitEqual
		}
		if label == "SplitUnequal" {
			return SplitUnequal
		}
		if label == "NotSplit" {
			return NotSplit
		}
	}

	return Undefined
}
