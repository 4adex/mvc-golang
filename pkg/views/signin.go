package views

import (
	"html/template"
)
func Signin() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/SignIn.html"))
	return temp
}