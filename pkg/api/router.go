package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/4adex/mvc-golang/pkg/middlewares"
	"github.com/4adex/mvc-golang/pkg/api/routes"
)

func Start() {
	r := mux.NewRouter()

	//Adding static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	//Public routes
	routes.RegisterPublicRoutes(r)

	//Authentication routes
	routes.RegisterAuthRoutes(r)
	r.Use(middleware.AuthMiddleware)

	//Admin Routes
	routes.RegisterAdminRoutes(r)

	http.ListenAndServe(":8000", r)
}