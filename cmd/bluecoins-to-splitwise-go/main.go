package main

import (
	"bluecoins-to-splitwise-go/pkg/bluecoins"
	"bluecoins-to-splitwise-go/pkg/model"
	"bluecoins-to-splitwise-go/pkg/splitwise"
	"fmt"
	"log"
	"math"
	"time"
)

func main() {

	bluecoinsService, err := bluecoins.NewBluecoinsService()
	if err != nil {
		log.Fatalf("Error creating Bluecoins service: %v", err)
	}
	transactions, err := bluecoinsService.GetTransactions()
	if err != nil {
		log.Fatalf("Error getting transactions: %v", err)
	}

	for _, t := range transactions {
		fmt.Printf("Transaction: %v\n", t)
		var addToSplitwise string
		fmt.Print("Add to Splitwise? (y/n): ")
		fmt.Scanln(&addToSplitwise)
		if addToSplitwise != "y" {
			continue
		}

		common := model.SplitwiseCommon{
			Description:     t.Description,
			Cost:            fmt.Sprintf("%f", math.Abs((float64(t.Amount)))),
			Date:            t.Date.Format(time.RFC3339),
			Currency_code:   "INR",
			Category_id:     12, // TODO: set the category as grocery
			Repeat_interval: model.Never,
			Details:         "Details",
		}
		// Add to Splitwise
		var isSplitEqual string
		fmt.Print("Is equal split? (y/n): ")
		fmt.Scanln(&isSplitEqual)
		var expense splitwise.SplitwiseService
		if isSplitEqual == "y" {
			// Create equal split expense
			expense = splitwise.ExpenseEqualGroupSplit{
				Group_id:        55886296,
				Split_equally:   true,
				SplitwiseCommon: common,
			}
		} else {
			// Create expense by shares
			expense = splitwise.ExpenseByShares{
				Group_id:             55886296,
				Users__0__user_id:    30164323,
				Users__0__paid_share: fmt.Sprintf("%f", t.Amount),
				Users__0__owed_share: "0",
				Users__1__first_name: "Soumya",
				Users__1__email:      "soumyam890@gmail.com",
				Users__1__paid_share: "0",
				Users__1__owed_share: fmt.Sprintf("%f", t.Amount),
				SplitwiseCommon:      common,
			}
		}

		err = expense.CreateExpense()

		if err != nil {
			log.Fatalf("Error creating expense: %v", err)
			fmt.Errorf("Error creating expense: %v", err)
		}
	}

	println("finished!")
}
