package views

import (
	"html/template"
)

func AdminRequest() *template.Template {
    temp := template.Must(template.ParseFiles("static/templates/ViewAdminRequests.html"))
    return temp
}