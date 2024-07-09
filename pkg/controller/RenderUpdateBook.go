package controller

import (
	"net/http"
	"strconv"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
)

//For viewing the page for updating the book
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