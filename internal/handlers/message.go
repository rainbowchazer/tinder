package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tinder-go/internal/database"
)

func GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	other := r.URL.Query().Get("user_email")
	if other == "" {
		http.Error(w, "Email пользователя обязателен", http.StatusBadRequest)
		return
	}

	messages, err := database.GetMessagesBetweenUsers(email, other)
	if err != nil {
		http.Error(w, "Ошибка при получении сообщений", http.StatusInternalServerError)
		log.Printf("Ошибка при получении: %v", err)
		return
	}
	json.NewEncoder(w).Encode(messages)
}
