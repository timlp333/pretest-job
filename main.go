package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	numCustomers    = 1000
	numTransactions = 5000
)

func main() {
	// Database connection
	dsn := "test:test@tcp(35.201.221.66:3306)/pretest"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Insert customers
	for i := 1; i <= numCustomers; i++ {
		name := fmt.Sprintf("Customer%d", i)
		email := fmt.Sprintf("customer%d@example.com", i)
		registrationDate := randomDate(time.Now().AddDate(-2, 0, 0), time.Now().AddDate(-1, 0, 0))

		_, err := db.Exec("INSERT INTO customers ( name, email, registration_date) VALUES (?, ?, ?)", name, email, registrationDate)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Map to keep track of the last transaction date for each customer
	customerLastTransactionDate := make(map[int]time.Time)
	customerTransactionCount := make(map[int]int)

	// Insert transactions
	for i := 1; i <= numTransactions; i++ {
		customerID := rand.Intn(numCustomers) + 1
		var registrationDateStr string
		err := db.QueryRow("SELECT registration_date FROM customers WHERE id = ?", customerID).Scan(&registrationDateStr)
		if err != nil {
			log.Fatal(err)
		}

		registrationDate, err := time.Parse("2006-01-02", registrationDateStr)
		if err != nil {
			log.Fatal(err)
		}

		// Determine the transaction date
		var transactionDate time.Time
		if lastDate, exists := customerLastTransactionDate[customerID]; exists {
			transactionDate = randomDate(lastDate, time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC))
		} else {
			transactionDate = randomDate(registrationDate, time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC))
		}

		amount := rand.Intn(1000) + 1
		customerTransactionCount[customerID]++
		transactionNumber := customerTransactionCount[customerID]

		_, err = db.Exec("INSERT INTO transactions ( customer_id, transaction_date, amount, transaction_number) VALUES ( ?, ?, ?, ?)", customerID, transactionDate, amount, transactionNumber)
		if err != nil {
			log.Fatal(err)
		}

		// Update the last transaction date for the customer
		customerLastTransactionDate[customerID] = transactionDate
	}
}

// randomDate generates a random date between start and end
func randomDate(start, end time.Time) time.Time {
	delta := end.Sub(start)
	sec := rand.Int63n(int64(delta.Seconds()))
	return start.Add(time.Duration(sec) * time.Second)
}
