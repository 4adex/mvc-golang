package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/4adex/mvc-golang/pkg/jwtutils"
	// "github.com/4adex/mvc-golang/pkg/messages"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/types"
	"github.com/4adex/mvc-golang/pkg/views"
	"golang.org/x/crypto/bcrypt"
)

func RenderSignin(w http.ResponseWriter, r *http.Request) {
	t := views.Signin()
	t.Execute(w, nil)
}

func RenderSignup(w http.ResponseWriter, r *http.Request) {
	t := views.Signup()
	t.Execute(w, nil)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  " ",
		MaxAge: -1,
	})
	
	jsonResponse(w, http.StatusOK, "/signin", "Logout Successful", "success")
	
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		
		jsonResponse(w, http.StatusInternalServerError, "/signin", "Unable to parse form successfully", "error")
		
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := models.GetUser(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		
		jsonResponse(w, http.StatusNotFound, "/", "User not found", "error")
		
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error comparing password: %v", err)
		
		jsonResponse(w, http.StatusUnauthorized, "/signin", "Invalid Password", "error")
		
		return
	}

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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Unable to parse form successfully", "error")
		return
	}

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")
	user.RequestStatus = "not_requested"

	
	isEmpty, err := models.IsUsersTableEmpty()
	if err != nil {
		log.Printf("Error checking users table: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Internal Server Error", "error")
		return
	}

	if isEmpty {
		user.Role = "admin"
	} else {
		user.Role = "client"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if (err != nil) {
		log.Printf("Error hashing password: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Internal Server Error", "error")
		return
	}

	user.Password = string(hashedPassword)
	err = models.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Internal Server Error", "error")
		return
	}

	token, err := jwtutils.GenerateJWT(user.Username, user.Email, user.Role, strconv.Itoa(user.ID))
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Internal Server Error", "error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	jsonResponse(w, http.StatusOK, "/", "Profile Created Successfully", "success")
}
