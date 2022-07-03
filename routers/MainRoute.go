package routers

import (
	"tes-mitramas/controllers"

	"github.com/gorilla/mux"
)

func MainRoute(route *mux.Router) {
	subroutes := route.PathPrefix("/api").Subrouter()
	subroutes.HandleFunc("/register", controllers.RegisterCont).Methods("POST")
	subroutes.HandleFunc("/login", controllers.LoginCont).Methods("POST")
	subroutes.HandleFunc("/logout", controllers.LogoutCont).Methods("POST")
	subroutes.HandleFunc("/checkin", controllers.CheckInCont).Methods("POST")
	subroutes.HandleFunc("/checkout", controllers.CheckOutCont).Methods("POST")
	subroutes.HandleFunc("/gettoken", controllers.GetTokenCont).Methods("GET")
}
