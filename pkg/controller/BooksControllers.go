package controller

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
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

//For viewing books on admin side
func RenderBooksAdmin(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetBooks()
	var username string = r.Context().Value("username").(string)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error loading from database", "error")
		return
	}
	
	data := make(map[string]interface{})
	data["Username"] = username
	data["Books"] = books

	t := views.BooksListAdmin()
	t.Execute(w, data)
}

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

func RenderAddBook(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	data := map[string]interface{}{
		"Username": username,
	}
	t := views.AddBook()
	t.Execute(w, data)
}

func HandleAddBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/addbook", "Error parsing form", "error")
		return
	}

	title := r.Form.Get("title")
	author := r.Form.Get("author")
	isbn := r.Form.Get("isbn")
	publicationYear := r.Form.Get("publication_year")
	totalCopies := r.Form.Get("total_copies")

	if title == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Title is required", "error")
		return
	}

	if author == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Author is required", "error")
		return
	}

	if isbn == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "ISBN is required", "error")
		return
	} else if len(isbn) > 13 {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "ISBN entered is too long", "error")
		return
	}

	if publicationYear == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Publication year is required", "error")
		return
	} else {
		year, err := strconv.Atoi(publicationYear)
		if err != nil || year < 1000 || year > time.Now().Year() {
			jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Invalid publication year", "error")
			return
		}
	}

	if totalCopies == "" {
		jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Total copies are required", "error")
		return
	} else {
		copies, err := strconv.Atoi(totalCopies)
		if err != nil || copies <= 0 {
			jsonResponse(w, http.StatusBadRequest, "/admin/addbook", "Total copies must be a positive number", "error")
			return
		}
	}

	err = models.InsertBook(title, author, isbn, publicationYear, totalCopies)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin/addbook", "Internal Server Error (Error Adding Book)", "error")
		return
	}
	jsonResponse(w, http.StatusCreated, "/admin/addbook", "Book added successfully", "success")
}