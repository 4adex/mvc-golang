package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/4adex/mvc-golang/pkg/jwtutils"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Unable to parse form successfully", "error")
		return
	}

	// Getting field values from the form data
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")
	user.RequestStatus = "not_requested"

	// Checking if any users exist before or not, which decides the role as user or admin
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

	// Generating hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Internal Server Error", "error")
		return
	}
	user.Password = string(hashedPassword)

	// Checking if username or email already exists
	exists, err := models.DoesUserExist(user.Username, user.Email)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Internal Server Error", "error")
		return
	}
	if exists {
		jsonResponse(w, http.StatusBadRequest, "/signup", "Username or Email already exists", "error")
		return
	}

	// Storing the details in the db
	err = models.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		jsonResponse(w, http.StatusInternalServerError, "/signup", "Internal Server Error", "error")
		return
	}

	// Generating the JWT token and storing it in cookie
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
