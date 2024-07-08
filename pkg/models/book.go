package models

import (
	"fmt"
	// "log"
	// "time"
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

func InsertBook(title, author, isbn, publicationYear, totalCopies string) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
		INSERT INTO books (title, author, isbn, publication_year, total_copies, available_copies)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(query, title, author, isbn, publicationYear, totalCopies, totalCopies)
	if err != nil {
		return err
	}

	return nil
}