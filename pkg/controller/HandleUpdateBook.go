package controller

import (
	"database/sql"
	"net/http"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/gorilla/mux"
)

//For making the request to update the book
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

	//error handling and checking for bad requests
	if title == "" || author == "" || isbn == "" || publicationYear == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks", "Parsed data is incomplete", "error")
		return
	} else if len(isbn) > 13 {
		jsonResponse(w, http.StatusBadRequest, "/admin/viewbooks", "ISBN entered is too long", "error")
		return
	}

	//updating book in db after ensuring data is clean
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