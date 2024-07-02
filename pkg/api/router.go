package api

import (
	"net/http"

	"github.com/4adex/mvc-golang/pkg/controller"
	"github.com/gorilla/mux"
	"github.com/4adex/mvc-golang/pkg/middlewares"
)

func Start() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/", controller.RenderHome).Methods("GET")
	r.HandleFunc("/viewbooks", controller.HandleViewBooks).Methods("GET")
	r.HandleFunc("/checkout/{id}", controller.CheckoutHandler).Methods("POST")

	r.HandleFunc("/signin", controller.RenderSignin).Methods("GET")
	r.HandleFunc("/signin", controller.SignInHandler).Methods("POST")
	r.HandleFunc("/signup",controller.RenderSignup).Methods("GET")
	r.HandleFunc("/signup", controller.SignUpHandler).Methods("POST")
	r.HandleFunc("/logout", controller.HandleLogout).Methods("POST")
	r.Use(middleware.AuthMiddleware)

	//first create list books route

	//separate router for admin routes and its middleware with it
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware)
    adminRouter.Use(middleware.AdminMiddleware)



	http.ListenAndServe(":8000", r)
}