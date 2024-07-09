package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/views"
)


func RenderSignup(w http.ResponseWriter, r *http.Request) {
	t := views.Signup()
	t.Execute(w, nil)
}