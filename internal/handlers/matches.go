package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tinder-go/internal/database"
)

func GetPotentialMatches(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email обязателен", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(context.Background(),
		`SELECT username, age, photo,email
		FROM users u
		LEFT JOIN likes l ON u.email = l.liked_email AND l.user_email = $1
		WHERE l.liked_email IS NULL AND u.email != $1
		LIMIT 10`, email)

	if err != nil {
		http.Error(w, "Ошибка поиска кандидатов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var candidates []UserProfile
	for rows.Next() {
		var user UserProfile
		err := rows.Scan(&user.Username, &user.Age, &user.Photo, &user.Email)
		if err != nil {
			log.Println("Ошибка при сканировании строки", err)
			continue
		}
		candidates = append(candidates, user)
	}

	json.NewEncoder(w).Encode(candidates)
}

func CheckForMatches(userEmail, likedEmail string) {
	var count int
	err := database.DB.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM likes
		WHERE user_email = $1 AND liked_email = $2 AND likes = true
		AND liked_email IN (SELECT user_email FROM likes WHERE user_email = $2 AND liked_email = $1 AND likes = true)`,
		userEmail, likedEmail).Scan(&count)

	if err != nil {
		log.Println("Ошибка при проверке мэтча:", err)
		return
	}

	if count > 0 {
		_, err := database.DB.Exec(context.Background(),
			`INSERT INTO matches (user1_email,user2_email) VALUES ($1,$2)`, userEmail, likedEmail)
		if err != nil {
			log.Println("Ошибка при сохранении мэтча:", err)
		} else {
			log.Printf("Новый мэтч: %s ❤️  %s", userEmail, likedEmail)
		}
	}
}

func LikeUsers(w http.ResponseWriter, r *http.Request) {
	var likeData struct {
		UserEmail  string `json:"user_email"`
		LikedEmail string `json:"liked_email"`
		Like       bool   `json:"like"`
	}

	if err := json.NewDecoder(r.Body).Decode(&likeData); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(context.Background(),
		`INSERT INTO likes (user_email, liked_email, likes)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_email,liked_email) DO UPDATE SET likes = $3`,
		likeData.UserEmail, likeData.LikedEmail, likeData.Like)

	if err != nil {
		http.Error(w, "Ошибка при сохранении лайка", http.StatusInternalServerError)
		return
	}

	CheckForMatches(likeData.UserEmail, likeData.LikedEmail)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Лайк/дизлайк сохранен")
}
