package views

import (
	"html/template"
)


func HomePage() *template.Template {
	temp := template.Must(template.ParseFiles("templates/home.html"))
	return temp
}

func BooksListClient() *template.Template {
	temp := template.Must(template.ParseFiles("templates/BookListClient.html"))
	return temp
}

func Signin() *template.Template {
	temp := template.Must(template.ParseFiles("templates/SignIn.html"))
	return temp
}

func Signup() *template.Template {
	temp := template.Must(template.ParseFiles("templates/SignUp.html"))
	return temp
}