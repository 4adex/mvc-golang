package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/views"
)

func RenderSignin(w http.ResponseWriter, r *http.Request) {
	t := views.Signin()
	t.Execute(w, nil)
}