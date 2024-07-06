package controller

import (
	"database/sql"
	// "github.com/4adex/mvc-golang/pkg/messages"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func RenderAdminHome(w http.ResponseWriter, r *http.Request) {
	var username string = r.Context().Value("username").(string)
	data := make(map[string]interface{})
	data["Username"] = username
	t := views.AdminHome()
	t.Execute(w, data)
}

func RenderBooksAdmin(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetBooks()
	var username string = r.Context().Value("username").(string)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error loading from database", "error")
		return
	}

	data := make(map[string]interface{})
	data["Username"] = username
	data["Books"] = books

	t := views.BooksListAdmin()
	t.Execute(w, data)
}

func RenderUpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks", "Invalid book ID", "error")
		return
	}

	book, err := models.GetBookByID(bookID)
	if err != nil {
		if err.Error() == "book not found" {
			jsonResponse(w, http.StatusNotFound, "/admin/viewbooks", "Book not found", "error")
			return
		}

		jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks", "Internal Server Error", "error")
		return
	}

	username := r.Context().Value("username").(string)
	data := make(map[string]interface{})
	data["Username"] = username
	data["Book"] = book
	t := views.UpdateBook()
	t.Execute(w, data)
}

func HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	err := r.ParseForm()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks", "Unable to parse form successfully", "error")
		return
	}
	title := r.FormValue("title")
	author := r.FormValue("author")
	isbn := r.FormValue("isbn")
	publicationYear := r.FormValue("publication_year")

	if title == "" || author == "" || isbn == "" || publicationYear == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks", "Parsed data is incomplete", "error")
		return
	} else if len(isbn) > 13 {
		jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks", "ISBN entered is too long", "error")
		return
	}

	err = models.UpdateBook(bookID, title, author, isbn, publicationYear)
	if err != nil {
		if err == sql.ErrNoRows {
			jsonResponse(w, http.StatusNotFound, "/admin/viewbooks", "Book Not Found", "error")
		} else {
			jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks", "Internal Server error", "error")
		}
		return
	}
	jsonResponse(w, http.StatusOK, "/admin/viewbooks", "Book updated successfully", "success")
}

func HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	isCheckedOut, err := models.IsBookCheckedOut(bookID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks", "Internal Server Error", "error")
		return
	}

	if isCheckedOut {
		jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks", "Cannot delete book that is currently checked out", "error")
		return
	}

	err = models.DeleteTransactionsByBookID(bookID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks", "Internal Server Error", "error")
		return
	}

	err = models.DeleteBookByID(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			jsonResponse(w, http.StatusNotFound, "/admin/viewbooks", "Book not found", "error")
		} else {
			jsonResponse(w, http.StatusInternalServerError, "/admin/viewbooks", "Internal Server Error", "error")
		}
		return
	}

	jsonResponse(w, http.StatusOK, "/admin/viewbooks", "Book deleted successfully", "success")
}

func RenderViewRequests(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)
	transactions, err := models.GetPendingTransactions(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin", "Internal Server Error", "error")
		return
	}

	data := map[string]interface{}{
		"Username":     username,
		"Transactions": transactions,
	}

	t := views.ViewRequests()
	t.Execute(w, data)
}

func HandleTransactionAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]
	action := vars["action"]

	if action != "accept" && action != "reject" {
		jsonResponse(w, http.StatusBadRequest, "/admin/viewrequests", "Invalid Action", "error")
		return
	}

	transaction, err := models.GetTransactionByID(transactionID)
	if err != nil {
		if err == sql.ErrNoRows {

			jsonResponse(w, http.StatusNotFound, "/admin/viewrequests", "Transaction not found", "error")
			return
		} else {

			jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests", "Internal Server error", "error")
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
			jsonResponse(w, http.StatusBadRequest, "/admin/viewrequests", "Not a valid action to do on transaction", "error")
			return
		}
	case "reject":
		if transaction.Status == "checkout_requested" {
			newStatus = "checkout_rejected"
		} else if transaction.Status == "checkin_requested" {
			newStatus = "checkin_rejected"
		} else {
			jsonResponse(w, http.StatusBadRequest, "/admin/viewrequests", "Not a valid action to do on transaction", "error")
			return
		}
	}

	if updateQuery != "" {
		_, err := models.UpdateBookAvailability(transaction.BookID, updateQuery)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests", "Internal Server error", "error")
			return
		}
	}

	err = models.UpdateTransactionStatusAdmin(transactionID, newStatus)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests", "Error Updating Transaction", "error")
		return
	}

	jsonResponse(w, http.StatusOK, "/admin/viewrequests", "Transaction updated successfully", "success")
}

func RenderAddBook(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	data := map[string]interface{}{
		"Username": username,
	}
	t := views.AddBook()
	t.Execute(w, data)
}

func HandleAddBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/addbook", "Error parsing form", "error")
		return
	}

	title := r.Form.Get("title")
	author := r.Form.Get("author")
	isbn := r.Form.Get("isbn")
	publicationYear := r.Form.Get("publication_year")
	totalCopies := r.Form.Get("total_copies")

	if title == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Title is required", "error")
		return
	}

	if author == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Author is required", "error")
		return
	}

	if isbn == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "ISBN is required", "error")
		return
	} else if len(isbn) > 13 {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "ISBN entered is too long", "error")
		return
	}

	if publicationYear == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Publication year is required", "error")
		return
	} else {
		year, err := strconv.Atoi(publicationYear)
		if err != nil || year < 1000 || year > time.Now().Year() {
			jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Invalid publication year", "error")
			return
		}
	}

	if totalCopies == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Total copies are required", "error")
		return
	} else {
		copies, err := strconv.Atoi(totalCopies)
		if err != nil || copies <= 0 {
			jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Total copies must be a positive number", "error")
			return
		}
	}

	err = models.InsertBook(title, author, isbn, publicationYear, totalCopies)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/addbook", "Internal Server Error (Error Adding Book)", "error")
		return
	}
	jsonResponse(w, http.StatusCreated, "/admin/addbook", "Book added successfully", "success")
}

func RenderAdminRequests(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, username, request_status
		FROM users
		WHERE role = 'client' AND request_status ="pending";
	`

	rows, err := models.GetPendingRequests(query)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin", "Internal Server Error", "error")
		return
	}

	username := r.Context().Value("username").(string)
	data := map[string]interface{}{
		"Requests": rows,
		"Username": username,
	}
	t := views.AdminRequest()
	t.Execute(w, data)
}

func HandleAcceptAdminRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	err := models.UpdateUserRoleAndStatus(userID, "admin", "accepted")
	if err != nil {
		if err == sql.ErrNoRows {
			jsonResponse(w, http.StatusNotFound, "/admin/adminrequests", "User not found", "error")
			return
		}
		jsonResponse(w, http.StatusInternalServerError, "/admin/adminrequests", "Internal Server Error", "error")
		return
	}
	jsonResponse(w, http.StatusOK, "/admin/adminrequests", "Admin access granted successfully", "success")
}

func HandleRejectAdminRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	err := models.UpdateUserRoleAndStatus(userID, "client", "rejected")
	if err != nil {
		if err == sql.ErrNoRows {
			// messages.SetFlash(w, r, "User not found", "error")
			jsonResponse(w, http.StatusNotFound, "/admin/adminrequests", "User not found", "error")
			return
		}

		jsonResponse(w, http.StatusInternalServerError, "/admin/adminrequests", "Internal Server Error", "error")

		return
	}

	jsonResponse(w, http.StatusOK, "/admin/adminrequests", "Admin access rejected successfully", "success")
}
