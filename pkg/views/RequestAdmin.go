package views

import (
	"html/template"
)

func ReqAdmin() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/reqAdmin.html"))
	return temp
}