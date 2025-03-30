package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tinder-go/internal/database"
)

type UserProfile struct {
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Age      int            `json:"age"`
	Photo    sql.NullString `json:"photo"`
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email обязателен", http.StatusBadRequest)
		return
	}

	var profile UserProfile
	err := database.DB.QueryRow(context.Background(),
		"SELECT username, email, age, photo FROM users WHERE email=$1", email).
		Scan(&profile.Username, &profile.Email, &profile.Age, &profile.Photo)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	if !profile.Photo.Valid {
		profile.Photo.String = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email обязателен", http.StatusBadRequest)
		return
	}

	var currentProfile UserProfile
	err := database.DB.QueryRow(context.Background(),
		"SELECT username, email, age, photo FROM users WHERE email=$1", email).
		Scan(&currentProfile.Username, &currentProfile.Email, &currentProfile.Age, &currentProfile.Photo)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка при получении текущих данных", http.StatusInternalServerError)
		return
	}

	var updateData struct {
		Username *string `json:"username"`
		Age      *int    `json:"age"`
		Photo    *string `json:"photo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	queryParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updateData.Username != nil {
		queryParts = append(queryParts, fmt.Sprintf("username=$%d", argIndex))
		args = append(args, *updateData.Username)
		argIndex++
	}

	if updateData.Age != nil {
		queryParts = append(queryParts, fmt.Sprintf("age=$%d", argIndex))
		args = append(args, *updateData.Age)
		argIndex++
	}

	if updateData.Photo != nil {
		queryParts = append(queryParts, fmt.Sprintf("photo=$%d", argIndex))
		args = append(args, *updateData.Photo)
		argIndex++
	} else if currentProfile.Photo.Valid {
		queryParts = append(queryParts, fmt.Sprintf("photo=$%d", argIndex))
		args = append(args, currentProfile.Photo.String)
		argIndex++
	}

	if len(queryParts) == 0 {
		http.Error(w, "Нет данных для обновления", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE email=$%d", strings.Join(queryParts, ", "), argIndex)
	args = append(args, email)

	_, err = database.DB.Exec(context.Background(), query, args...)

	if err != nil {
		http.Error(w, "Ошибка при обновлении профиля", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Профиль обновлен")
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email обязателен", http.StatusBadRequest)
		return
	}

	var id int

	err := database.DB.QueryRow(context.Background(),
		"SELECT id FROM users WHERE email=$1", email).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec(context.Background(),
		"DELETE FROM users WHERE id=$1", id)

	if err != nil {
		http.Error(w, "Ошибка удаления данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Профиль удален")
}
