package model

type RepeatInterval string

const (
	Never       RepeatInterval = "never"
	Weekly      RepeatInterval = "weekly"
	Fortnightly RepeatInterval = "fortnightly"
	Monthly     RepeatInterval = "monthly"
	Yearly      RepeatInterval = "yearly"
)

type SplitwiseCommon struct {
	Cost            string         `json:"cost"`            //A string representation of a decimal value, limited to 2 decimal places
	Description     string         `json:"description"`     // A short description of the expense
	Details         string         `json:"details"`         // Also known as notes
	Date            string         `json:"date"`            //The date and time the expense took place. May differ from `created_at`, example: 2012-05-02T13:00:00Z
	Repeat_interval RepeatInterval `json:"repeat_interval"` // repeat intervals : ["daily", "weekly", "monthly", "yearly", "never"]
	Currency_code   string         `json:"currency_code"`   //A currency code. Must be in the list from `get_currencies`
	Category_id     int            `json:"category_id"`     //A category id from `get_categories`
}
