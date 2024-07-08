package views

import (
	"html/template"
)

func AddBook() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/AdminAddBook.html"))
	return temp
}