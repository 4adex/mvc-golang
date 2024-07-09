package controller

import (
	"github.com/4adex/mvc-golang/pkg/models"
	"net/http"
	"github.com/gorilla/mux"
)


//For handling checkout requests from client side
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	bookId := mux.Vars(r)["id"]
	err := models.CreateCheckout(userId, bookId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error creating checkout", "error")
		return
	}
	jsonResponse(w, http.StatusOK, "/", "Checkout Requested Successfully", "success")
}