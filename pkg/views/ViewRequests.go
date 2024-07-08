package views

import (
	"html/template"
)

func ViewRequests() *template.Template {
    temp := template.Must(template.ParseFiles("static/templates/ViewCCRequests.html"))
    return temp
}