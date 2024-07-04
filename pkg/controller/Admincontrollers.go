package controller

import (
	// "database/sql"
	// "encoding/json"
	// "fmt"
	// "log"
	"net/http"
	// "time"

	"github.com/4adex/mvc-golang/pkg/messages"
	// "github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/views"
	// "github.com/gorilla/mux"
)

func RenderAdminHome(w http.ResponseWriter, r *http.Request) {
	var username string = r.Context().Value("username").(string)
	msg, msgType := messages.GetFlash(w, r)
	data := make(map[string]interface{})
	data["Username"] = username
	data["msg"] = msg
	data["msgType"] = msgType
	t := views.AdminHome()
	t.Execute(w, data)
}


