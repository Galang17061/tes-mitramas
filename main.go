package main

import (
	"fmt"
	"net/http"
	"os"
	"tes-mitramas/app"
	"tes-mitramas/models"
	"tes-mitramas/routers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func init() {
	models.BaseInit()
}

func main() {
	router := mux.NewRouter()
	models.DeleteTokenAndCheckout()

	routers.MainRoute(router)
	routers.AbsensiRoute(router)
	routers.AktivitasRoute(router)

	port := os.Getenv("PORT_APP")
	router.Use(app.JwtAuthentication)

	if port == "" {
		port = "9000"
	}

	optionsCode := handlers.OptionStatusCode(204)
	headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "X-Requested-With", "Content-Type", "Authorization", "X-CSRF-Token", "Content-Length", "Accept-Encoding", "Accept"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})

	fmt.Println(port)

	err := http.ListenAndServe(os.Getenv("HOST_APP")+":"+port, handlers.CORS(optionsCode, originsOk, headersOk, methodsOk)(router))
	if err != nil {
		fmt.Print(err)
	}
}
