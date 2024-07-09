package controller

import (
	"net/http"
)


func HandleLogout(w http.ResponseWriter, r *http.Request) {

	//Deleting the cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  " ",
		MaxAge: -1,
	})
	
	//sending the response and redirect to /signin
	jsonResponse(w, http.StatusOK, "/signin", "Logout Successful", "success")
	
}