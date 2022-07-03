package routers

import (
	"tes-mitramas/controllers"

	"github.com/gorilla/mux"
)

func AktivitasRoute(route *mux.Router) {
	subroutes := route.PathPrefix("/api").Subrouter()
	subroutes.HandleFunc("/act/create", controllers.CreateAktivitas).Methods("POST")
	subroutes.HandleFunc("/act/edit", controllers.EditAktivitas).Methods("POST")
	subroutes.HandleFunc("/act/delete", controllers.DeleteAktivitas).Methods("POST")
}
