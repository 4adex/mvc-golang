package views

import (
	"html/template"
)

func UpdateBook() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/AdminUpdateBook.html"))
	return temp
}