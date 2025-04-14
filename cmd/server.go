package main

import (
	"log"
	"net/http"
	"tinder-go/internal/database"
	"tinder-go/internal/handlers"
	"tinder-go/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.ConnectDB()

	r := mux.NewRouter()

	r.Handle("/register", middleware.ValidateRegister(http.HandlerFunc(handlers.RegisterUser))).Methods("POST")
	r.Handle("/login", middleware.ValidateLogin(http.HandlerFunc(handlers.LoginUser))).Methods("POST")
	r.Handle("/messages", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetMessageHistory))).Methods("GET")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/profile", handlers.GetProfile).Methods("GET")
	protected.HandleFunc("/profile/update", handlers.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/profile/delete", handlers.DeleteProfile).Methods("DELETE")
	protected.HandleFunc("/matches", handlers.GetPotentialMatchesWithParam).Methods("GET")
	protected.HandleFunc("/like", handlers.LikeUsers).Methods("POST")
	protected.HandleFunc("/profile/interests", handlers.UpdateInterest).Methods("PUT")

	http.HandleFunc("/ws", handlers.ChatHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders:   []string{"Autharization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(r)

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", handler)
}
