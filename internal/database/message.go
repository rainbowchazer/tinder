package database

import (
	"context"
	"time"
	"tinder-go/internal/models"
)

func SaveMessage(senderEmail, receiverEmail, message string) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO messages (sender_email, receiver_email, message, timestamp)
	VALUES ($1,$2,$3,$4)`, senderEmail, receiverEmail, message, time.Now())
	return err
}

func GetMessagesBetweenUsers(email1, email2 string) ([]models.Message, error) {
	rows, err := DB.Query(context.Background(), `
	SELECT id,sender_email,receiver_email,message,timestamp
	FROM messages
	WHERE (sender_email=$1 AND receiver_email=$2)
	OR (sender_email=$2 AND receiver_email =$1)
	ORDER BY timestamp
	`, email1, email2)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.SenderEmail, &m.ReceiverEmail, &m.Message, &m.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
