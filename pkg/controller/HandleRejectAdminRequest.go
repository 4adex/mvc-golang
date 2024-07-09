package controller

import (
	"database/sql"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)


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
