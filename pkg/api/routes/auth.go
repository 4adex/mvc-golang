package routes

import (
    "github.com/gorilla/mux"
    "github.com/4adex/mvc-golang/pkg/controller"
)

func RegisterAuthRoutes(r *mux.Router) {
    r.HandleFunc("/signin", controller.RenderSignin).Methods("GET")
    r.HandleFunc("/signin", controller.SignInHandler).Methods("POST")
    r.HandleFunc("/signup", controller.RenderSignup).Methods("GET")
    r.HandleFunc("/signup", controller.SignUpHandler).Methods("POST")
    r.HandleFunc("/logout", controller.HandleLogout).Methods("POST")
}