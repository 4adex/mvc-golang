package controller

import (
	"database/sql"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)

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