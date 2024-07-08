package views

import (
	"html/template"
)

func HoldingsClient() *template.Template {
    temp := template.Must(template.ParseFiles("static/templates/viewholdings.html"))
    return temp
}