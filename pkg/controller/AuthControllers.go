package controller

import (
	// "fmt"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"github.com/4adex/mvc-golang/pkg/jwtutils"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/types"
	"github.com/4adex/mvc-golang/pkg/views"
	"golang.org/x/crypto/bcrypt"
	// "github.com/4adex/mvc-golang/pkg/views"
)

func RenderSignin(w http.ResponseWriter, r *http.Request) {
	// Implementation for rendering the sign-in page
	t := views.Signin()
	t.Execute(w, nil)
}

func RenderSignup(w http.ResponseWriter, r *http.Request) {
	// Implementation for rendering the sign-in page
	t := views.Signup()
	t.Execute(w, nil)
}

func HandleLogout(w http.ResponseWriter, r *http.Request){
	//Clearing the cookies does the job right, and maybe redirection is also needed
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: "",
		MaxAge: -1,
	})
	res := types.Response{
		Message: "Logout Successful",
		Type: "success",
	}
	json.NewEncoder(w).Encode(res)

}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := models.GetUser(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error comparing password: %v", err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := jwtutils.GenerateJWT(user.Username, user.Email, user.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}


func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")
	user.Role = "client"
	user.RequestStatus = "not_requested"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	err = models.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := jwtutils.GenerateJWT(user.Username, user.Email, user.Role)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}