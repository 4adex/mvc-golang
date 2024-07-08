package views

import (
	"html/template"
)

func AdminHome() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/AdminHome.html"))
	return temp
}