package controller

import (
	"net/http"
	"strconv"
	"time"
	"github.com/4adex/mvc-golang/pkg/models"
)


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