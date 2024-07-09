package controller

import (
	"database/sql"
	"github.com/4adex/mvc-golang/pkg/models"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)



//For handling checkin requests from client side
func HandleCheckin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]
	transaction, err := models.GetTransactionByID(transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			jsonResponse(w, http.StatusNotFound, "/", "Transaction not found", "error")
			return
		} else {
			jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server error", "error")
			return
		}
	}

	if transaction.Status == "checkin_requested" {
		jsonResponse(w, http.StatusBadRequest, "/", "Checkin is already requested for this transaction", "error")
		return
	} else if transaction.Status == "checkin_accepted" {
		jsonResponse(w, http.StatusBadRequest, "/", "Checkin is already accepted for this transaction", "error")
		return
	} else if transaction.Status != "checkout_accepted" {
		jsonResponse(w, http.StatusBadRequest, "/", "Transaction must be checked out first", "error")
		return
	}
	err = models.UpdateTransactionStatus(transactionID, "checkin_requested", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
		return
	}
	jsonResponse(w, http.StatusOK, "/", "Checkin requested successfully", "success")
}