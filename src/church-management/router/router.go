package router

import (
	"net/http"

	"church-management/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/auth/change-password", handlers.ChangePassword).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/dashboard", handlers.GetDashboardData).Methods("GET")
	return r
}

func StartServer() {
	r := InitRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
}
