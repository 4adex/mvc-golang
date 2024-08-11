package controller

import (
	"database/sql"
	"fmt"
	"strconv"

	// "fmt"
	"net/http"

	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
)

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

	switch action {
	case "accept":
		if transaction.Status == "checkout_requested" {
			bookIDInt, err := strconv.Atoi(transaction.BookID)
			if err != nil {
				jsonResponse(w, http.StatusBadRequest, "/", "Invalid book ID", "error")
				return
			}
			availableCopies, _, err := models.GetBookCopies(bookIDInt)
			if err != nil {
				jsonResponse(w, http.StatusInternalServerError, "/", "Error retrieving book copies", "error")
				return
			}
			if availableCopies<1 {
				jsonResponse(w, http.StatusBadRequest, "/admin/viewrequests", "No available copies to checkout", "error")
				return
			}
			newStatus = "checkout_accepted"
			_, err = models.DecreaseBookQuantity(transaction.BookID)
			if err != nil {
				jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests", "Error updating book quantity", "error")
				return
			}
		} else if transaction.Status == "checkin_requested" {
			newStatus = "returned"
			_, err := models.IncreaseBookQuantity(transaction.BookID)
			fmt.Println("boook is   ",transaction.BookID)
			fmt.Println(transaction)
			if err != nil {
				jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests", "Error updating book quantity", "error")
				return
			}
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

	err = models.UpdateTransactionStatusAdmin(transactionID, newStatus)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/viewrequests", "Error Updating Transaction", "error")
		return
	}

	jsonResponse(w, http.StatusOK, "/admin/viewrequests", "Transaction updated successfully", "success")
}