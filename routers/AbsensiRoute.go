package routers

import (
	"tes-mitramas/controllers"

	"github.com/gorilla/mux"
)

func AbsensiRoute(route *mux.Router) {
	subroutes := route.PathPrefix("/api").Subrouter()
	subroutes.HandleFunc("/absen/list", controllers.RiwayatAbsensi).Methods("GET")
	subroutes.HandleFunc("/absen/listsbytanggal", controllers.RiwayatAbsensiByTanggal).Methods("GET")
}
