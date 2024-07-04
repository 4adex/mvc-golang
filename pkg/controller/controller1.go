package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/4adex/mvc-golang/pkg/messages"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
)

func jsonResponse(w http.ResponseWriter, status int, redirect string) {
	response := map[string]string{
		"redirect": redirect,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func RenderHome(w http.ResponseWriter, r *http.Request) {
	var username string = r.Context().Value("username").(string)
	msg, msgType := messages.GetFlash(w, r)
	fmt.Println("Messages are", msg, msgType)
	data := make(map[string]interface{})
	data["Username"] = username
	data["msg"] = msg
	data["msgType"] = msgType
	t := views.HomePage()
	t.Execute(w, data)
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
    userId := r.Context().Value("id").(string)
    bookId := mux.Vars(r)["id"]
    fmt.Println(userId, bookId)

    err := models.CreateCheckout(userId, bookId)
    if err != nil {
        messages.SetFlash(w, r, "Error Creating Checkout", "error")
        jsonResponse(w, http.StatusInternalServerError, "/viewbooks")
        return
    }

    messages.SetFlash(w, r, "Checkout Requested Successfully", "success")
    jsonResponse(w, http.StatusOK, "/viewbooks")
}

func HandleViewBooks(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(data)

	t := views.BooksListClient()
	t.Execute(w, data)
}

func HandleViewHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)

	histories, err := models.GetHistory(userID)
	if err != nil {
		messages.SetFlash(w, r, "Error loading from database", "error")
		jsonResponse(w, http.StatusInternalServerError, "/")
		return
	}
	msg, msgType := messages.GetFlash(w, r)
	data := make(map[string]interface{})
	data["Username"] = username
	data["Transactions"] = histories
	data["msg"] = msg
	data["msgType"] = msgType

	t := views.HistoryListClient()
	t.Execute(w, data)
}

func HandleViewHoldings(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)
	holdings, err := models.GetHoldings(userId)
	if err != nil {
		messages.SetFlash(w, r, "Internal Server Error", "error")
		jsonResponse(w, http.StatusInternalServerError, "/")
		return
	}
	msg, msgType := messages.GetFlash(w, r)
	data := map[string]interface{}{
		"Username":     username,
		"Transactions": holdings,
	}
	data["msg"] = msg
	data["msgType"] = msgType

	t := views.HoldingsClient()
	t.Execute(w, data)
}

func HandleCheckin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	transaction, err := models.GetTransactionByID(transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			messages.SetFlash(w, r, "Transaction not found", "error")
			jsonResponse(w, http.StatusNotFound, "/")
			return
		} else {
			log.Println(err)
			messages.SetFlash(w, r, "Internal Server error", "error")
			jsonResponse(w, http.StatusInternalServerError, "/")
			return
		}
	}

	if transaction.Status == "checkin_requested" {
		messages.SetFlash(w, r, "Checkin is already requested for this transaction", "error")
		jsonResponse(w, http.StatusBadRequest, "/")
		return
	} else if transaction.Status == "checkin_accepted" {
		messages.SetFlash(w, r, "Checkin is already accepted for this transaction", "error")
		jsonResponse(w, http.StatusBadRequest, "/")
		return
	} else if transaction.Status != "checkout_accepted" {
		messages.SetFlash(w, r, "Transaction must be checked out first", "error")
		jsonResponse(w, http.StatusBadRequest, "/")
		return
	}

	err = models.UpdateTransactionStatus(transactionID, "checkin_requested", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println(err)
		messages.SetFlash(w, r, "Internal Server Error", "error")
		jsonResponse(w, http.StatusInternalServerError, "/")
		return
	}
	messages.SetFlash(w, r, "Checkin requested successfully", "success")
	jsonResponse(w, http.StatusOK, "/")
}

func HandleAdminRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)

	role, requestStatus, err := models.GetUserRequestStatus(userID)
	if err != nil {
		messages.SetFlash(w, r, "Internal Server Error", "error")
		jsonResponse(w, http.StatusInternalServerError, "/")
		return
	}

	if role == "admin" {
		messages.SetFlash(w, r, "You are already an admin", "error")
		jsonResponse(w, http.StatusBadRequest, "/")
		return
	}

	if requestStatus == "rejected" || requestStatus == "not_requested" {
		err := models.UpdateUserRequestStatus(userID, "pending")
		if err != nil {
			messages.SetFlash(w, r, "Internal Server Error", "error")
			jsonResponse(w, http.StatusInternalServerError, "/")
			return
		}
		messages.SetFlash(w, r, "Admin request sent successfully", "success")
		jsonResponse(w, http.StatusOK, "/")
		return
	}

	messages.SetFlash(w, r, "Admin request already exists or has been accepted", "error")
	jsonResponse(w, http.StatusBadRequest, "/")
}

func RenderReqAdmin(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)
	_, requestStatus, err := models.GetUserRequestStatus(userID)
	if err != nil {
		messages.SetFlash(w, r, "Internal Server Error", "error")
		jsonResponse(w, http.StatusInternalServerError, "/")
		return
	}
	if requestStatus == "rejected" || requestStatus == "not_requested" {
		msg, msgType := messages.GetFlash(w, r)
		data := make(map[string]interface{})
		data["Username"] = username
		data["msg"] = msg
		data["msgType"] = msgType
		t := views.ReqAdmin()
		t.Execute(w, data)
		return
	}
	messages.SetFlash(w, r, "Admin request already exists or has been accepted", "error")
	jsonResponse(w, http.StatusBadRequest, "/")
}
