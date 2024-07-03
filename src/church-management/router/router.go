package router

import (
	"net/http"

	"church-management/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	return r
}

func StartServer() {
	r := InitRouter()
	handler := cors.Default().Handler(r)
	http.ListenAndServe(":8080", handler)
}
