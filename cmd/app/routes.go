package main

import (
	"github.com/aarsinh/auth-golang/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", app.login).Methods("POST")
	r.HandleFunc("/signup", app.signup).Methods("POST")
	r.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(app.accessProtected))).Methods("GET")
	return r
}
