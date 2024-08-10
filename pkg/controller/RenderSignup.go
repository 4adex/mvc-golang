package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/views"
)


func RenderSignup(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("token")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t := views.Signup()
	t.Execute(w, nil)
}