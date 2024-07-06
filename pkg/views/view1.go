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

func HistoryListClient() *template.Template {
	temp := template.Must(template.ParseFiles("templates/ViewHistory.html"))
	return temp
}

func HoldingsClient() *template.Template {
    temp := template.Must(template.ParseFiles("templates/viewholdings.html"))
    return temp
}

func ReqAdmin() *template.Template {
	temp := template.Must(template.ParseFiles("templates/reqAdmin.html"))
	return temp
}

func AdminHome() *template.Template {
	temp := template.Must(template.ParseFiles("templates/AdminHome.html"))
	return temp
}

func BooksListAdmin() *template.Template {
	temp := template.Must(template.ParseFiles("templates/BookListAdmin.html"))
	return temp
}

func UpdateBook() *template.Template {
	temp := template.Must(template.ParseFiles("templates/AdminUpdateBook.html"))
	return temp
}

func AddBook() *template.Template {
	temp := template.Must(template.ParseFiles("templates/AdminAddBook.html"))
	return temp
}

func ViewRequests() *template.Template {
    temp := template.Must(template.ParseFiles("templates/ViewCCRequests.html"))
    return temp
}

func AdminRequest() *template.Template {
    temp := template.Must(template.ParseFiles("templates/ViewAdminRequests.html"))
    return temp
}





