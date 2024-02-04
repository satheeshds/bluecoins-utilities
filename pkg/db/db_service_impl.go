package db

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Implementing the DBService interface
type DBServiceImpl struct {
	db *sql.DB
}

func NewDBService(db string) (*DBServiceImpl, error) {
	dbConn, err := sql.Open("sqlite3", db)
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

func (m *DBServiceImpl) GetTransactions(after time.Time) ([]model.Transaction, error) {
	// Implement your logic here
	// For now, we'll just return an empty slice and nil error
	query := `SELECT tt.transactionstableid, tt.date, tt.amount, tt.categoryid, it.itemname 
			FROM transactionstable tt inner join itemtable it on it.itemtableid = tt.itemid 
			where tt.date between ? and ?;`
	rows, err := m.db.Query(query, after, time.Now().AddDate(0, 0, 1))
	if err != nil {
		return nil, err
	}

	var transactions []model.Transaction
	var amount int
	for rows.Next() {
		var transaction model.Transaction
		err = rows.Scan(&transaction.ID, &transaction.Date, &amount, &transaction.Category, &transaction.Description)
		if err != nil {
			return nil, err
		}
		transaction.Amount = float32(amount) / 1000000

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
