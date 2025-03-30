package handlers

import (
	"context"
	"log"
	"net/http"
	"tinder-go/internal/database"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender   string `json:"sender_email"`
	Receiver string `json:"receiver_email"`
	Content  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при установке WebSocket-соединения:", err)
		return
	}
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Ошибка при чтении сообщения:", err)
			break
		}

		_, err = database.DB.Exec(context.Background(),
			"INSERT INTO messages (sender_email, receiver_email, message) VALUES ($1,$2,$3)",
			msg.Sender, msg.Receiver, msg.Content)
		if err != nil {
			log.Println("Ошибка сохранения в БД:", err)
			continue
		}

		err = conn.WriteJSON(map[string]string{"status": "Сообщение отправлено"})
		if err != nil {
			log.Println("Ошибка при отправке ответа", err)
			break
		}
	}
}
