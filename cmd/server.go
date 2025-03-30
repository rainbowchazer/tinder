package main

import (
	"log"
	"net/http"
	"tinder-go/internal/database"
	"tinder-go/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/profile", handlers.GetProfile).Methods("GET")
	r.HandleFunc("/profile/update", handlers.UpdateProfile).Methods("PUT")
	r.HandleFunc("/profile/delete", handlers.DeleteProfile).Methods("DELETE")
	r.HandleFunc("/matches", handlers.GetPotentialMatches).Methods("GET")
	r.HandleFunc("/like", handlers.LikeUsers).Methods("POST")

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", r)
}
