package routes

import (
    "github.com/gorilla/mux"
    "github.com/4adex/mvc-golang/pkg/controller"
    "github.com/4adex/mvc-golang/pkg/middlewares"
)

func RegisterAdminRoutes(r *mux.Router) {
    adminRouter := r.PathPrefix("/admin").Subrouter()

    adminRouter.HandleFunc("/dashboard", controller.RenderAdminHome).Methods("GET")
    adminRouter.HandleFunc("/viewbooks", controller.RenderBooksAdmin).Methods("GET")
    adminRouter.HandleFunc("/update/{id}", controller.RenderUpdateBook).Methods("GET")
    adminRouter.HandleFunc("/update/{id}", controller.HandleUpdateBook).Methods("POST")
    adminRouter.HandleFunc("/delete/{id}", controller.HandleDeleteBook).Methods("POST")
    adminRouter.HandleFunc("/viewrequests", controller.RenderViewRequests).Methods("GET")
    adminRouter.HandleFunc("/transaction/{id}/{action}", controller.HandleTransactionAction).Methods("POST")
    adminRouter.HandleFunc("/addbook", controller.RenderAddBook).Methods("GET")
    adminRouter.HandleFunc("/addbook", controller.HandleAddBook).Methods("POST")
    adminRouter.HandleFunc("/adminrequests", controller.RenderAdminRequests).Methods("GET")
    adminRouter.HandleFunc("/adminrequest/accept/{id}", controller.HandleAcceptAdminRequest).Methods("POST")
    adminRouter.HandleFunc("/adminrequest/reject/{id}", controller.HandleRejectAdminRequest).Methods("POST")

    adminRouter.Use(middleware.AuthMiddleware)
    adminRouter.Use(middleware.AdminMiddleware)
}