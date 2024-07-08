package routes

import (
    "github.com/gorilla/mux"
    "github.com/4adex/mvc-golang/pkg/controller"
)

func RegisterPublicRoutes(r *mux.Router) {
    r.HandleFunc("/", controller.RenderHome).Methods("GET")
    r.HandleFunc("/viewbooks", controller.HandleViewBooks).Methods("GET")
    r.HandleFunc("/checkout/{id}", controller.CheckoutHandler).Methods("POST")
    r.HandleFunc("/history", controller.HandleViewHistory).Methods("GET")
    r.HandleFunc("/viewholdings", controller.HandleViewHoldings).Methods("GET")
    r.HandleFunc("/checkin/{id}", controller.HandleCheckin).Methods("POST")
    r.HandleFunc("/reqadmin", controller.RenderReqAdmin).Methods("GET")
    r.HandleFunc("/reqadmin", controller.HandleAdminRequest).Methods("POST")
}