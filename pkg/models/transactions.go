package models

import (
	"fmt"
	"log"
	"time"
	"database/sql"
	"github.com/4adex/mvc-golang/pkg/types"
)


func CreateCheckout(userID string, bookID string) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	currentDatetime := time.Now().Format("2006-01-02 15:04:05")
	query := "INSERT INTO `transactions`(`user_id`, `book_id`, `status`,`checkout_time`) VALUES (?, ?, ?, ?)"
	_, err = db.Query(query, userID, bookID, "checkout_requested", currentDatetime)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetHistory(userID string) ([]types.History, error) {
	db, err := Connection()
	if err != nil {
		fmt.Printf("error %s connecting to Database", err)
		return nil, err
	}
	defer db.Close()

	query := `
    SELECT t.transaction_id, b.title, t.status, t.checkout_time, t.checkin_time
    FROM transactions t
    JOIN books b ON t.book_id = b.id
    WHERE t.user_id = ?`

	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Printf("error %s querying the database", err)
		return nil, err
	}
	defer rows.Close()

	var histories []types.History
	for rows.Next() {
		var history types.History
		var checkinTime sql.NullString

		err := rows.Scan(&history.TransactionID, &history.Title, &history.Status, &history.CheckoutTime, &checkinTime)
		if err != nil {
			log.Printf("error %s scanning the row", err)
			continue
		}

		if checkinTime.Valid {
			history.CheckinTime = checkinTime.String
		} else {
			history.CheckinTime = ""
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error %s iterating over rows", err)
		return nil, err
	}

	fmt.Println(histories)

	return histories, nil
}

func GetHoldings(userID string) ([]types.Holding, error) {
    db, err := Connection()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    query := `
    SELECT t.transaction_id, b.title, b.author, t.checkout_time
    FROM books b
    JOIN transactions t ON b.id = t.book_id
    WHERE t.status IN ('checkout_accepted', 'checkin_requested') AND t.user_id = ?;`
    
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var holdings []types.Holding
    for rows.Next() {
        var holding types.Holding
        err := rows.Scan(&holding.TransactionID, &holding.Title, &holding.Author, &holding.CheckoutTime)
        if err != nil {
            return nil, err
        }
        holdings = append(holdings, holding)
    }

    return holdings, nil
}

func GetTransactionByID(transactionID string) (types.Transaction, error) {
    db, err := Connection()
    if err != nil {
        return types.Transaction{}, err
    }
    defer db.Close()

    query := "SELECT * FROM `transactions` WHERE `transaction_id` = ?"
    row := db.QueryRow(query, transactionID)

    var transaction types.Transaction
    err = row.Scan(&transaction.TransactionID, &transaction.BookID, &transaction.UserID, &transaction.Status, &transaction.CheckoutTime, &transaction.CheckinTime)
    if err != nil {
        return types.Transaction{}, err
    }

    return transaction, nil
}

func UpdateTransactionStatus(transactionID, status, checkinTime string) error {
    db, err := Connection()
    if err != nil {
        return err
    }
    defer db.Close()

    query := "UPDATE `transactions` SET `status` = ?, `checkin_time` = ? WHERE `transaction_id` = ?"
    _, err = db.Exec(query, status, checkinTime, transactionID)
    if err != nil {
        return err
    }

    return nil
}

func DeleteTransactionsByBookID(bookID string) error {
    db, err := Connection()
    if err != nil {
        return err
    }
    defer db.Close()

    query := "DELETE FROM transactions WHERE book_id = ?"
    _, err = db.Exec(query, bookID)
    return err
}

func GetPendingTransactions(userID string) ([]types.PendingTransaction, error) {
    db, err := Connection()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    query := `
      SELECT t.transaction_id, b.title, t.status, t.checkout_time, t.checkin_time
      FROM transactions t
      JOIN books b ON t.book_id = b.id
      WHERE t.status NOT IN ('checkout_rejected', 'checkin_rejected', 'returned', 'checkout_accepted')
      AND t.user_id = ?
    `

    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []types.PendingTransaction
    for rows.Next() {
        var transaction types.PendingTransaction
        err := rows.Scan(&transaction.TransactionID, &transaction.Title, &transaction.Status, &transaction.CheckoutTime, &transaction.CheckinTime)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }

    return transactions, nil
}

func UpdateTransactionStatusAdmin(transactionID string, status string) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE transactions SET status = ? WHERE transaction_id = ?"
	_, err = db.Exec(query, status, transactionID)
	if err != nil {
		return err
	}

	return nil
}

func GetPendingRequests() ([]types.AdminRequest, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

    query := `
		SELECT id, username, request_status
		FROM users
		WHERE role = 'client' AND request_status ="pending";
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []types.AdminRequest
	for rows.Next() {
		var request types.AdminRequest
		err := rows.Scan(&request.ID, &request.Username, &request.RequestStatus)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}