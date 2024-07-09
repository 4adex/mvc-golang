package controller

import (
	"net/http"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
)

func HandleViewHoldings(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	username := r.Context().Value("username").(string)
	holdings, err := models.GetHoldings(userId)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "/", "Internal Server Error", "error")
		return
	}
	data := map[string]interface{}{
		"Username":     username,
		"Transactions": holdings,
	}
	t := views.HoldingsClient()
	t.Execute(w, data)
}


