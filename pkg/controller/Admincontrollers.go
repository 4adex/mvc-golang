package controller

import (
	// "database/sql"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	// "time"
	"strconv"
	"github.com/4adex/mvc-golang/pkg/messages"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
	"database/sql"
)

func RenderAdminHome(w http.ResponseWriter, r *http.Request) {
	var username string = r.Context().Value("username").(string)
	msg, msgType := messages.GetFlash(w, r)
	data := make(map[string]interface{})
	data["Username"] = username
	data["msg"] = msg
	data["msgType"] = msgType
	t := views.AdminHome()
	t.Execute(w, data)
}

func RenderBooksAdmin(w http.ResponseWriter, r *http.Request){
	books, err := models.GetBooks()
	var username string = r.Context().Value("username").(string)
	if err != nil {
		messages.SetFlash(w, r, "Error loading from database", "error")
		jsonResponse(w, http.StatusInternalServerError, "/")
		return
	}

	data := make(map[string]interface{})
	msg, msgType := messages.GetFlash(w, r)
	data["Username"] = username
	data["Books"] = books
	data["msg"] = msg
	data["msgType"] = msgType

	t := views.BooksListAdmin()
	t.Execute(w, data)
}

func RenderUpdateBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bookID, err := strconv.Atoi(vars["id"])
    if err != nil {
		messages.SetFlash(w, r, "Invalid book ID", "error")
        jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks")
        return
    }

    book, err := models.GetBookByID(bookID)
    if err != nil {
        if err.Error() == "book not found" {
			messages.SetFlash(w, r, "Book not found", "error")
            jsonResponse(w, http.StatusNotFound, "/admin/viewbooks")
            return
        }
        fmt.Println("Error fetching book:", err)
		messages.SetFlash(w, r, "Internal Server Error", "error")
        jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks")
        return
    }

    username := r.Context().Value("username").(string)
	msg, msgType := messages.GetFlash(w,r)
	data := make(map[string]interface{})
	data["Username"] = username
	data["msg"] = msg
	data["msgType"] = msgType
	data["Book"] = book

	t := views.UpdateBook()
	t.Execute(w,data)
}


func HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bookID := vars["id"]

    err := r.ParseForm()
    if err != nil {
        log.Printf("Error parsing form: %v", err)
        messages.SetFlash(w, r, "Unable to parse form successfully", "error")
        jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks")
        return
    }

    title := r.FormValue("title")
    author := r.FormValue("author")
    isbn := r.FormValue("isbn")
    publicationYear := r.FormValue("publication_year")

    if title == "" || author == "" || isbn == "" || publicationYear == "" {
        messages.SetFlash(w, r, "Parsed data is incomplete", "error")
        jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks")
        return
    } else if len(isbn) > 13 {
        messages.SetFlash(w, r, "ISBN entered is too long", "error")
        jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks")
        return
    }

    err = models.UpdateBook(bookID, title, author, isbn, publicationYear)
    if err != nil {
        if err == sql.ErrNoRows {
            messages.SetFlash(w, r, "Book Not Found", "error")
            jsonResponse(w, http.StatusNotFound, "/admin/viewbooks")
        } else {
            log.Printf("Error updating book: %v", err)
            messages.SetFlash(w, r, "Internal Server error", "error")
            jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks")
        }
        return
    }

    messages.SetFlash(w, r, "Book updated successfully", "success")
    jsonResponse(w, http.StatusOK, "/admin/viewbooks")
}

func HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bookID := vars["id"]

    // Check if the book is currently checked out
    isCheckedOut, err := models.IsBookCheckedOut(bookID)
    if err != nil {
        log.Printf("Error checking book status: %v", err)
        messages.SetFlash(w, r, "Internal Server Error", "error")
        jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks")
        return
    }

    if isCheckedOut {
        messages.SetFlash(w, r, "Cannot delete book that is currently checked out", "error")
        jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks")
        return
    }

    // Delete transactions related to the book
    err = models.DeleteTransactionsByBookID(bookID)
    if err != nil {
        log.Printf("Error deleting transactions: %v", err)
        messages.SetFlash(w, r, "Internal Server Error", "error")
        jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks")
        return
    }

    // Delete the book
    err = models.DeleteBookByID(bookID)
    if err != nil {
        if err == sql.ErrNoRows {
            messages.SetFlash(w, r, "Book not found", "error")
            jsonResponse(w, http.StatusNotFound, "/admin/viewbooks")
        } else {
            log.Printf("Error deleting book: %v", err)
            messages.SetFlash(w, r, "Internal Server Error", "error")
            jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks")
        }
        return
    }

    messages.SetFlash(w, r, "Book deleted successfully", "success")
    jsonResponse(w, http.StatusOK, "/admin/viewbooks")
}


func RenderViewRequests(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("id").(string)
    username := r.Context().Value("username").(string)

    transactions, err := models.GetPendingTransactions(userID)
    if err != nil {
        log.Printf("Error fetching pending transactions: %v", err)
        messages.SetFlash(w, r, "Internal Server Error", "error")
        jsonResponse(w, http.StatusInternalServerError, "/admin")
        return
    }

    msg, msgType := messages.GetFlash(w, r)

    data := map[string]interface{}{
        "Username":     username,
        "Transactions": transactions,
        "msg":          msg,
        "msgType":      msgType,
    }

    t := views.ViewRequests()
    t.Execute(w, data)
}

func HandleTransactionAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]
	action := vars["action"]

	if action != "accept" && action != "reject" {
		messages.SetFlash(w, r, "Invalid Action", "error")
		jsonResponse(w, http.StatusBadRequest, "/admin/viewrequests")
		return
	}

	transaction, err := models.GetTransactionByID(transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			messages.SetFlash(w, r, "Transaction not found", "error")
			jsonResponse(w, http.StatusNotFound, "/admin/viewrequests")
			return
		} else {
			log.Println(err)
			messages.SetFlash(w, r, "Internal Server error", "error")
			jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests")
			return
		}
	}

	var newStatus string
	var updateQuery string

	switch action {
	case "accept":
		if transaction.Status == "checkout_requested" {
			newStatus = "checkout_accepted"
			updateQuery = "UPDATE books SET available_copies = available_copies - 1 WHERE id = ?"
		} else if transaction.Status == "checkin_requested" {
			newStatus = "returned"
			updateQuery = "UPDATE books SET available_copies = available_copies + 1 WHERE id = ?"
		} else {
			messages.SetFlash(w, r, "Not a valid action to do on transaction", "error")
			jsonResponse(w, http.StatusBadRequest, "/admin/viewrequests")
			return
		}
	case "reject":
		if transaction.Status == "checkout_requested" {
			newStatus = "checkout_rejected"
		} else if transaction.Status == "checkin_requested" {
			newStatus = "checkin_rejected"
		} else {
			messages.SetFlash(w, r, "Not a valid action to do on transaction", "error")
			jsonResponse(w, http.StatusBadRequest, "/admin/viewrequests")
			return
		}
	}

	if updateQuery != "" {
		_, err := models.UpdateBookAvailability(transaction.BookID, updateQuery)
		if err != nil {
			log.Println(err)
			messages.SetFlash(w, r, "Internal Server error", "error")
			jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests")
			return
		}
	}

	err = models.UpdateTransactionStatusAdmin(transactionID, newStatus)
	if err != nil {
		log.Println(err)
		messages.SetFlash(w, r, "Error Updating Transaction", "error")
		jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests")
		return
	}

	messages.SetFlash(w, r, "Transaction updated successfully", "success")
	jsonResponse(w, http.StatusOK, "/admin/viewrequests")
}


func RenderAddBook(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	msg, msgType := messages.GetFlash(w, r)

	data := map[string]interface{}{
		"Username": username,
		"msg":      msg,
		"msgType":  msgType,
	}

	t := views.AddBook()
	t.Execute(w, data)
}
