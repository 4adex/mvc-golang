package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/4adex/mvc-golang/pkg/jwtutils"
	"github.com/4adex/mvc-golang/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signin", "Unable to parse form successfully", "error")
		return
	}

	//Getting username and password fields from form
	username := r.FormValue("username")
	password := r.FormValue("password")

	//Retrieving data from the db
	user, err := models.GetUser(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		jsonResponse(w, http.StatusNotFound, "/", "User not found", "error")
		return
	}

	//comparing hash with password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error comparing password: %v", err)
		jsonResponse(w, http.StatusUnauthorized, "/signin", "Invalid Password", "error")
		return
	}

	//generating JWT token and storing it in cookie
	token, err := jwtutils.GenerateJWT(user.Username, user.Email, user.Role, strconv.Itoa(user.ID))
	if err != nil {
		log.Printf("Error generating token: %v", err)	
		jsonResponse(w, http.StatusInternalServerError, "/signin", "Error Generating Token", "error")	
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	jsonResponse(w, http.StatusOK, "/", "Logged In Successfully", "success")
	
}