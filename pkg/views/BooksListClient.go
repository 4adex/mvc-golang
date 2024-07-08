package views

import (
	"html/template"
)

func BooksListClient() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/BookListClient.html"))
	return temp
}