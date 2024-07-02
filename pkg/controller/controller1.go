package controller

import (
	// "fmt"
	"net/http"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/types"
	"github.com/4adex/mvc-golang/pkg/views"
)


type PageData struct {
    Username string
    Books    []types.Book
}

func RenderHome(w http.ResponseWriter, r *http.Request) {
	t := views.HomePage()
	t.Execute(w, nil)
}

func HandleViewBooks(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetBooks()
	var username string = r.Context().Value("username").(string)
	if err != nil {
		http.Error(w, "Error loading from database", http.StatusInternalServerError)
		return
	}
	data := PageData{
		Username: username,
		Books:    books,
	}
	t := views.BooksListClient()
	t.Execute(w, data)
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for checkout logic
}

