package views

import (
	"html/template"
)

func Signup() *template.Template {
	temp := template.Must(template.ParseFiles("static/templates/SignUp.html"))
	return temp
}