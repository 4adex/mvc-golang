package controller

import (
	"database/sql"
	"net/http"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/gorilla/mux"
)

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