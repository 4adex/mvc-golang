package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/views"
)

func RenderHome(w http.ResponseWriter, r *http.Request) {
	var username string = r.Context().Value("username").(string)
	role := r.Context().Value("role").(string)
	data := make(map[string]interface{})
	data["Username"] = username
	data["Role"] = role
	t := views.HomePage()
	t.Execute(w, data)
}

func RenderAdminHome(w http.ResponseWriter, r *http.Request) {
	var username string = r.Context().Value("username").(string)
	data := make(map[string]interface{})
	data["Username"] = username
	t := views.AdminHome()
	t.Execute(w, data)
}