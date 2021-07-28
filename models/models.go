package models

import "time"

type Message struct {
	SenderUsername string
	Content        string
	CreatedAt      time.Time
}

type User struct {
	Username  string
	CreatedAt time.Time
}
