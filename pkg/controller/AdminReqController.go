package controller

import (
	"database/sql"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	"github.com/gorilla/mux"
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


//For accepting the admin request
func HandleAcceptAdminRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	err := models.UpdateUserRoleAndStatus(userID, "admin", "accepted")
	if err != nil {
		if err == sql.ErrNoRows {
			jsonResponse(w, http.StatusNotFound, "/admin/adminrequests", "User not found", "error")
			return
		}
		jsonResponse(w, http.StatusInternalServerError, "/admin/adminrequests", "Internal Server Error", "error")
		return
	}
	jsonResponse(w, http.StatusOK, "/admin/adminrequests", "Admin access granted successfully", "success")
}

//For rejecting the admin request
func HandleRejectAdminRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	err := models.UpdateUserRoleAndStatus(userID, "client", "rejected")
	if err != nil {
		if err == sql.ErrNoRows {
			// messages.SetFlash(w, r, "User not found", "error")
			jsonResponse(w, http.StatusNotFound, "/admin/adminrequests", "User not found", "error")
			return
		}
		jsonResponse(w, http.StatusInternalServerError, "/admin/adminrequests", "Internal Server Error", "error")
		return
	}

	jsonResponse(w, http.StatusOK, "/admin/adminrequests", "Admin access rejected successfully", "success")
}
