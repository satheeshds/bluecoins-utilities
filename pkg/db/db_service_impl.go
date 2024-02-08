package db

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"database/sql"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Implementing the DBService interface
type DBServiceImpl struct {
	db *sql.DB
}

func NewDBService(db string) (*DBServiceImpl, error) {
	dbConn, err := sql.Open("sqlite", db)
	if err != nil {
		return nil, err
	}
	dbService := &DBServiceImpl{}
	dbService.db = dbConn
	return dbService, nil
}

func (m *DBServiceImpl) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *DBServiceImpl) GetTransactions(after time.Time, accountId int) ([]model.Transaction, error) {
	// Implement your logic here
	// For now, we'll just return an empty slice and nil error
	query := `
		SELECT 
        	tt.transactionstableid, 
        	tt.date, 
        	tt.amount, 
        	tt.categoryid, 
        	it.itemname,
        	(SELECT GROUP_CONCAT(lt.labelname) FROM labelstable lt WHERE lt.transactionidlabels = tt.transactionstableid) as labels
    	FROM 
        	transactionstable tt 
    	INNER JOIN 
        	itemtable it 
    	ON 
        	it.itemtableid = tt.itemid 
    	WHERE 
        	accountid = ? 
    	AND 
        	tt.date between ? and ?
		ORDER BY tt.date ASC;
		`
	rows, err := m.db.Query(query, accountId, after, time.Now().AddDate(0, 0, 1))
	if err != nil {
		return nil, err
	}

	var transactions []model.Transaction
	var amount int
	var labels string
	for rows.Next() {
		var transaction model.Transaction
		err = rows.Scan(&transaction.ID, &transaction.Date, &amount, &transaction.Category, &transaction.Description, &labels)
		if err != nil {
			return nil, err
		}
		transaction.Amount = float32(amount) / 1000000
		transaction.Labels = strings.Split(labels, ",")

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *DBServiceImpl) GetAccounts() ([]model.Account, error) {
	// Implement your logic here
	// For now, we'll just return an empty slice and nil error
	query := `SELECT accountname, accountstableid FROM accountstable;`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}

	var accounts []model.Account
	for rows.Next() {
		var account model.Account
		err = rows.Scan(&account.Name, &account.ID)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
