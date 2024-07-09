package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
)

//For viewing books on client side
func HandleViewBooks(w http.ResponseWriter, r *http.Request) {

	//Getting books from the 
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