package views

import (
	"html/template"
)

func HistoryListClient() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/ViewHistory.html"))
	return temp
}