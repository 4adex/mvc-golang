package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
)


func HandleViewHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)

	histories, err := models.GetHistory(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Error loading from database", "error")
		return
	}
	data := make(map[string]interface{})
	data["Username"] = username
	data["Transactions"] = histories
	t := views.HistoryListClient()
	t.Execute(w, data)
}