package controller

import (
	// "fmt"
	"net/http"
	"strconv"

	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/gorilla/mux"
)

//For handling checkout requests from client side
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	bookId := mux.Vars(r)["id"]

	// Convert bookId to int
	bookIDInt, err := strconv.Atoi(bookId)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "/", "Invalid book ID", "error")
		return
	}

	// Check available copies before proceeding with checkout
	availableCopies, _, err := models.GetBookCopies(bookIDInt)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error retrieving book copies", "error")
		return
	}

	if availableCopies < 1 {
		jsonResponse(w, http.StatusBadRequest, "/", "No available copies for checkout", "error")
		return
	}

	err = models.CreateCheckout(userId, bookId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error creating checkout", "error")
		return
	}

	jsonResponse(w, http.StatusOK, "/", "Checkout Requested Successfully", "success")
}
