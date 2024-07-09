package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/views"
)

func RenderAddBook(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	data := map[string]interface{}{
		"Username": username,
	}
	t := views.AddBook()
	t.Execute(w, data)
}