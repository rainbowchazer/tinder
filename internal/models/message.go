package models

import "time"

type Message struct {
	ID            int       `json:"id"`
	SenderEmail   string    `json:"sender_email"`
	ReceiverEmail string    `json:"receiver_email"`
	Message       string    `json:"message"`
	Timestamp     time.Time `json:"timestamp"`
}
