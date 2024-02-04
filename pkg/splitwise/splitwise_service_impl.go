package splitwise

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type ExpenseEqualGroupSplit struct {
	model.SplitwiseCommon
	Group_id      int  `json:"group_id"`      // The group to put this expense in.
	Split_equally bool `json:"split_equally"` // split equally among all members of the group
}

type ExpenseByShares struct {
	model.SplitwiseCommon
	Group_id             int    `json:"group_id"` // The group to put this expense in.
	Users__0__user_id    int    `json:"users__0__user_id"`
	Users__0__paid_share string `json:"users__0__paid_share"` // "Decimal amount as a string with 2 decimal places. The amount this user paid for the expense"
	Users__0__owed_share string `json:"users__0__owed_share"` // "Decimal amount as a string with 2 decimal places. The amount this user owes for the expense"
	Users__1__first_name string `json:"users__1__first_name"`
	Users__1__last_name  string `json:"users__1__last_name"`
	Users__1__email      string `json:"users__1__email"`
	Users__1__paid_share string `json:"users__1__paid_share"` // "Decimal amount as a string with 2 decimal places. The amount this user paid for the expense"
	Users__1__owed_share string `json:"users__1__owed_share"` // "Decimal amount as a string with 2 decimal places. The amount this user owes for the expense"
}

func (e ExpenseByShares) CreateExpense() error {
	//convert expense to json
	jsonData, err := json.Marshal(e)
	if err != nil {
		return err
	}

	// make api call
	err = makeApiCall(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (e ExpenseEqualGroupSplit) CreateExpense() error {
	//convert expense to json
	jsonData, err := json.Marshal(e)
	if err != nil {
		return err
	}

	// make api call
	err = makeApiCall(jsonData)
	if err != nil {
		return err
	}

	return nil

}

func makeApiCall(data []byte) error {

	fmt.Println(string(data))

	// make api call
	req, err := http.NewRequest("POST", "https://www.splitwise.com/api/v3.0/create_expense", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	apiKey := os.Getenv("SPLITWISE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("SPLITWISE_API_KEY not set")
	}

	// set the auth token
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check the status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response %v", resp)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
