package routes

import (
    "github.com/gorilla/mux"
    "github.com/4adex/mvc-golang/pkg/controller"
)

func RegisterAdminRoutes(r *mux.Router) {
    r.HandleFunc("/admin/dashboard", controller.RenderAdminHome).Methods("GET")
    r.HandleFunc("/admin/viewbooks", controller.RenderBooksAdmin).Methods("GET")
    r.HandleFunc("/admin/update/{id}", controller.RenderUpdateBook).Methods("GET")
    r.HandleFunc("/admin/update/{id}", controller.HandleUpdateBook).Methods("POST")
    r.HandleFunc("/admin/delete/{id}", controller.HandleDeleteBook).Methods("POST")
    r.HandleFunc("/admin/viewrequests", controller.RenderViewRequests).Methods("GET")
    r.HandleFunc("/admin/transaction/{id}/{action}", controller.HandleTransactionAction).Methods("POST")
    r.HandleFunc("/admin/addbook", controller.RenderAddBook).Methods("GET")
    r.HandleFunc("/admin/addbook", controller.HandleAddBook).Methods("POST")
    r.HandleFunc("/admin/adminrequests", controller.RenderAdminRequests).Methods("GET")
    r.HandleFunc("/admin/adminrequest/accept/{id}", controller.HandleAcceptAdminRequest).Methods("POST")
    r.HandleFunc("/admin/adminrequest/reject/{id}", controller.HandleRejectAdminRequest).Methods("POST")
}