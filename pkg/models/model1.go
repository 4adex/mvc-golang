package models

import (
	// "database/sql"
	"fmt"
	"log"
	"time"
	"database/sql"

	"github.com/4adex/mvc-golang/pkg/types"
)

func GetBooks() ([]types.Book, error) {
	db, err := Connection()
	if err != nil {
		fmt.Printf("error %s connecting to Database", err)
		return nil, err
	}
	query := "SELECT * FROM books"
	rows, err := db.Query(query)
	db.Close()
	if err != nil {
		fmt.Printf("error %s querying the database", err)
		return nil, err
	}
	var BooksList []types.Book
	for rows.Next() {
		var book types.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.TotalCopies, &book.AvailableCopies)
		if err != nil {
			fmt.Printf("error %s scanning the row", err)
			continue
		}
		BooksList = append(BooksList, book)
	}
	return BooksList, nil
}

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

// Model
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

func GetUserRequestStatus(userID string) (string, string, error) {
    db, err := Connection()
    if err != nil {
        return "", "", err
    }
    defer db.Close()

    var role, requestStatus string
    query := "SELECT role, request_status FROM users WHERE id = ?"
    err = db.QueryRow(query, userID).Scan(&role, &requestStatus)
    if err != nil {
        return "", "", err
    }
    return role, requestStatus, nil
}

func UpdateUserRequestStatus(userID string, status string) error {
    db, err := Connection()
    if err != nil {
        return err
    }
    defer db.Close()

    query := "UPDATE users SET request_status = ? WHERE id = ?"
    _, err = db.Exec(query, status, userID)
    if err != nil {
        return err
    }
    return nil
}


func GetBookByID(bookID int) (types.Book, error) {
	var book types.Book
	db, err := Connection()
    if err != nil {
        return book,err
    }
    defer db.Close()
    query := "SELECT * FROM books WHERE id = ?"
    err = db.QueryRow(query, bookID).Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.TotalCopies, &book.AvailableCopies)
    if err != nil {
        if err == sql.ErrNoRows {
            return book, fmt.Errorf("book not found")
        }
        return book, err
    }
    return book, nil
}



func UpdateBook(bookID, title, author, isbn, publicationYear string) error {
    db, err := Connection()
    if err != nil {
        return err
    }
    defer db.Close()

    query := `
      UPDATE books
      SET title = ?, author = ?, isbn = ?, publication_year = ? 
      WHERE id = ?
    `
    result, err := db.Exec(query, title, author, isbn, publicationYear, bookID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}


func IsBookCheckedOut(bookID string) (bool, error) {
    db, err := Connection()
    if err != nil {
        return false, err
    }
    defer db.Close()

    query := "SELECT COUNT(*) AS count FROM transactions WHERE book_id = ? AND status = 'checkout_accepted'"
    var count int
    err = db.QueryRow(query, bookID).Scan(&count)
    if err != nil {
        return false, err
    }

    return count > 0, nil
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

func DeleteBookByID(bookID string) error {
    db, err := Connection()
    if err != nil {
        return err
    }
    defer db.Close()

    query := "DELETE FROM books WHERE id = ?"
    result, err := db.Exec(query, bookID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
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


func UpdateBookAvailability(bookID string, query string) (sql.Result, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	result, err := db.Exec(query, bookID)
	if err != nil {
		return nil, err
	}

	return result, nil
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
