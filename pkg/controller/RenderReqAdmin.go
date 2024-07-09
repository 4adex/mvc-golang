package controller

import (
	// "database/sql"
	"github.com/4adex/mvc-golang/pkg/models"
	// "log"
	"net/http"
	// "time"
	"github.com/4adex/mvc-golang/pkg/views"
	// "github.com/gorilla/mux"
)


//For viewing page of requesting admins
func RenderReqAdmin(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)
	_, requestStatus, err := models.GetUserRequestStatus(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
		return
	}
	if requestStatus == "rejected" || requestStatus == "not_requested" {
		data := make(map[string]interface{})
		data["Username"] = username
		t := views.ReqAdmin()
		t.Execute(w, data)
		return
	}
	jsonResponse(w, http.StatusBadRequest, "/", "Admin request already exists or has been accepted", "error")
}