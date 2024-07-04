package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/4adex/mvc-golang/pkg/jwtutils"
	"github.com/4adex/mvc-golang/pkg/messages"
	"github.com/4adex/mvc-golang/pkg/models"
	"github.com/4adex/mvc-golang/pkg/types"
	"github.com/4adex/mvc-golang/pkg/views"
	"golang.org/x/crypto/bcrypt"
)

func RenderSignin(w http.ResponseWriter, r *http.Request) {
	t := views.Signin()
	msg, msgType := messages.GetFlash(w, r)
	fmt.Println("Messages are", msg, msgType)
	data := make(map[string]interface{})
	data["msg"] = msg
	data["msgType"] = msgType
	t.Execute(w, data)
}

func RenderSignup(w http.ResponseWriter, r *http.Request) {
	t := views.Signup()
	msg, msgType := messages.GetFlash(w, r)
	fmt.Println("Messages are", msg, msgType)
	data := make(map[string]interface{})
	data["msg"] = msg
	data["msgType"] = msgType
	t.Execute(w, data)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  " ",
		MaxAge: -1,
	})
	messages.SetFlash(w, r, "Logout Successful", "success")
	jsonResponse(w, http.StatusOK, "/signin")
	// http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303 See Other
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		messages.SetFlash(w, r, "Unable to parse form successfully", "error")
		jsonResponse(w, http.StatusInternalServerError, "/signin")
		// http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303 See Other
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := models.GetUser(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		messages.SetFlash(w, r, "User not found", "error")
		jsonResponse(w, http.StatusNotFound, "/")
		// http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303 See Other
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error comparing password: %v", err)
		messages.SetFlash(w, r, "Invalid Password", "error")
		jsonResponse(w, http.StatusUnauthorized, "/signin")
		// http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303 See Other
		return
	}

	token, err := jwtutils.GenerateJWT(user.Username, user.Email, user.Role, strconv.Itoa(user.ID))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		messages.SetFlash(w, r, "Error Generating Token", "error")
		jsonResponse(w, http.StatusInternalServerError, "/signin")
		// http.Redirect(w, r, "/signin", http.StatusSeeOther) // 303 See Other
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	messages.SetFlash(w, r, "Logged In Successfully", "success")
	jsonResponse(w, http.StatusOK, "/")
	// http.Redirect(w, r, "/", http.StatusSeeOther) // 303 See Other
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		messages.SetFlash(w, r, "Unable to parse form successfully", "error")
		jsonResponse(w, http.StatusInternalServerError, "/signup")
		// http.Redirect(w, r, "/signup", http.StatusSeeOther) // 303 See Other
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
		messages.SetFlash(w, r, "Internal Server Error", "error")
		jsonResponse(w, http.StatusInternalServerError, "/signup")
		// http.Redirect(w, r, "/signup", http.StatusSeeOther) // 303 See Other
		return
	}

	user.Password = string(hashedPassword)
	err = models.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		messages.SetFlash(w, r, "Internal Server Error", "error")
		jsonResponse(w, http.StatusInternalServerError, "/signup")
		// http.Redirect(w, r, "/signup", http.StatusSeeOther) // 303 See Other
		return
	}

	token, err := jwtutils.GenerateJWT(user.Username, user.Email, user.Role, strconv.Itoa(user.ID))
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		messages.SetFlash(w, r, "Internal Server Error", "error")
		jsonResponse(w, http.StatusInternalServerError, "/signup")
		// http.Redirect(w, r, "/signup", http.StatusSeeOther) // 303 See Other
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})

	messages.SetFlash(w, r, "Profile Created Successfully", "success")
	jsonResponse(w, http.StatusOK, "/")
	// http.Redirect(w, r, "/", http.StatusSeeOther) // 303 See Other
}
