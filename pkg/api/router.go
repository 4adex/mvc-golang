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
	r.HandleFunc("/history", controller.HandleViewHistory).Methods("GET")
	r.HandleFunc("/viewholdings", controller.HandleViewHoldings).Methods("GET")
	r.HandleFunc("/checkin/{id}", controller.HandleCheckin).Methods("POST")
	r.HandleFunc("/reqadmin", controller.RenderReqAdmin).Methods("GET")
	r.HandleFunc("/reqadmin", controller.HandleAdminRequest).Methods("POST")

	r.HandleFunc("/signin", controller.RenderSignin).Methods("GET")
	r.HandleFunc("/signin", controller.SignInHandler).Methods("POST")
	r.HandleFunc("/signup",controller.RenderSignup).Methods("GET")
	r.HandleFunc("/signup", controller.SignUpHandler).Methods("POST")
	r.HandleFunc("/logout", controller.HandleLogout).Methods("POST")
	r.Use(middleware.AuthMiddleware) 
	// r.Use(middleware.FlashMiddleware)

	//first create list books route

	//separate router for admin routes and its middleware with it
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/dashboard", controller.RenderAdminHome).Methods("GET")
	adminRouter.HandleFunc("/viewbooks", controller.RenderBooksAdmin).Methods("GET")
	adminRouter.HandleFunc("/update/{id}",controller.RenderUpdateBook).Methods("GET")
	adminRouter.HandleFunc("/update/{id}",controller.HandleUpdateBook).Methods("POST")
	adminRouter.HandleFunc("/delete/{id}",controller.HandleDeleteBook).Methods("POST")
	adminRouter.HandleFunc("/viewrequests",controller.RenderViewRequests).Methods("GET")
	adminRouter.HandleFunc("/transaction/{id}/{action}",controller.HandleTransactionAction).Methods("POST")
	adminRouter.HandleFunc("/addbook", controller.RenderAddBook).Methods("GET")
	adminRouter.HandleFunc("/addbook", controller.HandleAddBook).Methods("POST")
	adminRouter.HandleFunc("/adminrequests",controller.RenderAdminRequests).Methods("GET")
	adminRouter.HandleFunc("/adminrequest/accept/{id}", controller.HandleAcceptAdminRequest).Methods("POST")
	adminRouter.HandleFunc("/adminrequest/reject/{id}", controller.HandleRejectAdminRequest).Methods("POST")





	adminRouter.Use(middleware.AuthMiddleware)
    adminRouter.Use(middleware.AdminMiddleware)



	http.ListenAndServe(":8000", r)
}