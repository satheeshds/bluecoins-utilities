package splitwise

import (
	"bluecoins-to-splitwise-go/pkg/gdrive"
	"bluecoins-to-splitwise-go/pkg/model"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type SplitwiseServiceImpl struct {
	Expense      Expense
	gdrive       gdrive.GDriveServiceImpl
	accountState map[int]time.Time
}

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

// type SplitwiseState struct {
// 	// AccountId       int       `json:"account_id"`
// 	LastExpenseDate time.Time `json:"last_expense_date"`
// }

func NewSplitwiseService() (*SplitwiseServiceImpl, error) {
	driveService, err := gdrive.NewGDriveService()
	if err != nil {
		fmt.Printf("error creating gdrive service: %v \n", err)
	}
	err = driveService.DownloadFile("splitwise.json", "splitwise.json")
	if err != nil {
		fmt.Printf("error downloading splitwise.json: %v\n", err)
	}

	file, err := os.Open("splitwise.json")
	if err != nil {
		fmt.Printf("error opening splitwise.json: %v\n", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	state := make(map[int]time.Time)
	err = decoder.Decode(&state)
	if err != nil {
		fmt.Printf("error decoding splitwise.json: %v \n", err)
	}

	service := &SplitwiseServiceImpl{}
	service.gdrive = *driveService
	service.accountState = state
	return service, nil
}

func (s *SplitwiseServiceImpl) GetLastExpenseDate(accountId int) (time.Time, error) {
	defaultDate := time.Now().Local().AddDate(0, 0, -7)
	date, ok := s.accountState[accountId]
	if ok {
		return date, nil
	}
	return defaultDate, nil
}

func (s *SplitwiseServiceImpl) SetLastExpenseDate(accountId int, date time.Time) {
	s.accountState[accountId] = date
}

func (e ExpenseByShares) Create() error {
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

func (e ExpenseEqualGroupSplit) Create() error {
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

func (s *SplitwiseServiceImpl) Close() {

	// Perform cleanup here.
	// For example, if you have a file that needs to be closed, you can close it here.
	// If there's an error, return it.
	jsonData, err := json.Marshal(s.accountState)
	if err != nil {
		log.Printf("error marshalling account state: %v", err)
	}

	file, err := os.Create("splitwise.json")
	if err != nil {
		log.Printf("error creating splitwise.json: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Printf("error writing to splitwise.json: %v", err)
	}

	err = s.gdrive.UploadFile("splitwise.json", "splitwise.json")
	if err != nil {
		log.Printf("error uploading splitwise.json: %v", err)
	}

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
