package views

import (
	"html/template"
)

func BooksListAdmin() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/BookListAdmin.html"))
	return temp
}
