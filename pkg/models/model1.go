package models

import (
	// "database/sql"
	"fmt"
	"time"

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
		return err
	}
	return nil
}


