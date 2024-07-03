package main

import (
	"church-management/router"
	"church-management/utils"
	"log"
	"net/http"

	"church-management/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	_, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")

	// Configuração do CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
	router.StartServer()
}
