package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
)

func jsonResponse(w http.ResponseWriter, status int, redirect string, flashMessage string, flashType string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"redirect":     redirect,
		"flashMessage": flashMessage,
		"flashType":    flashType,
	})
}

func RenderHome(w http.ResponseWriter, r *http.Request) {
	var username string = r.Context().Value("username").(string)
	role := r.Context().Value("role").(string)
	data := make(map[string]interface{})
	data["Username"] = username
	data["Role"] = role
	t := views.HomePage()
	t.Execute(w, data)
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	bookId := mux.Vars(r)["id"]
	fmt.Println(userId, bookId)

	err := models.CreateCheckout(userId, bookId)
	if err != nil {
		log.Printf("Error creating checkout: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/", "Error creating checkout", "error")
		return
	}

	jsonResponse(w, http.StatusOK, "/", "Checkout Requested Successfully", "success")
}

func HandleViewBooks(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetBooks()
	var username string = r.Context().Value("username").(string)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error loading from database", "error")
		return
	}

	data := make(map[string]interface{})
	data["Username"] = username
	data["Books"] = books

	t := views.BooksListClient()
	t.Execute(w, data)
}

func HandleViewHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)

	histories, err := models.GetHistory(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error loading from database", "error")
		return
	}

	data := make(map[string]interface{})
	data["Username"] = username
	data["Transactions"] = histories
	t := views.HistoryListClient()
	t.Execute(w, data)
}

func HandleViewHoldings(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)
	holdings, err := models.GetHoldings(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
		return
	}
	data := map[string]interface{}{
		"Username":     username,
		"Transactions": holdings,
	}

	t := views.HoldingsClient()
	t.Execute(w, data)
}

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

func HandleAdminRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)

	role, requestStatus, err := models.GetUserRequestStatus(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
		return
	}

	if role == "admin" {
		jsonResponse(w, http.StatusBadRequest, "/", "You are already an admin", "error")
		return
	}

	if requestStatus == "rejected" || requestStatus == "not_requested" {
		err := models.UpdateUserRequestStatus(userID, "pending")
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
			return
		}

		jsonResponse(w, http.StatusOK, "/", "Admin request sent successfully", "success")
		return
	}

	jsonResponse(w, http.StatusBadRequest, "/", "Admin request already exists or has been accepted", "error")
}

func RenderReqAdmin(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)
	_, requestStatus, err := models.GetUserRequestStatus(userID)
	if err != nil {

		jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
		return
	}
	if requestStatus == "rejected" || requestStatus == "not_requested" {
		data := make(map[string]interface{})
		data["Username"] = username
		t := views.ReqAdmin()
		t.Execute(w, data)
		return
	}

	jsonResponse(w, http.StatusBadRequest, "/", "Admin request already exists or has been accepted", "error")
}
