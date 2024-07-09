package controller

import (
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"net/http"
)

//For getting list of users asking for admin requests
func RenderAdminRequests(w http.ResponseWriter, r *http.Request) {

	//Getting pending request sfrom db
	rows, err := models.GetPendingRequests()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/admin", "Internal Server Error", "error")
		return
	}

	//parsing variables in view and executing them
	username := r.Context().Value("username").(string)
	data := map[string]interface{}{
		"Requests": rows,
		"Username": username,
	}
	t := views.AdminRequest()
	t.Execute(w, data)
}