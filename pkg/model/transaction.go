package model

import "time"

// Transaction represents a Bluecoins transaction
type Transaction struct {
	ID          string
	Date        time.Time
	Amount      float32
	Category    int
	Description string
}
