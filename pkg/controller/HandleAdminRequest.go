package controller

import (
	"github.com/4adex/mvc-golang/pkg/models"
	"net/http"
)


//For handling admin requests from client
func HandleAdminRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	role, requestStatus, err := models.GetUserRequestStatus(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
		return
	}
	if role == "admin" {
		jsonResponse(w, http.StatusBadRequest, "/", "You are already an admin", "error")
		return
	}
	if requestStatus == "rejected" || requestStatus == "not_requested" {
		err := models.UpdateUserRequestStatus(userID, "pending")
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
			return
		}
		jsonResponse(w, http.StatusOK, "/", "Admin request sent successfully", "success")
		return
	}
	jsonResponse(w, http.StatusBadRequest, "/", "Admin request already exists or has been accepted", "error")
}
