package routes

import (
	"github.com/gorilla/mux"
	"github.com/nishant01/procard-go-api/controllers"
	"net/http"
)

var (
	Router = mux.NewRouter()
)

func MapUrls() {
	Router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
	Router.HandleFunc("/register", controllers.AccountController.Register).Methods(http.MethodPost)
	Router.HandleFunc("/login", controllers.AccountController.Login).Methods(http.MethodPost)
}
