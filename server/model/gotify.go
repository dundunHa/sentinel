package model

import (
	"gorm.io/gorm"
)

type GotifyMessage struct {
	gorm.Model
	AppID           int `gorm:"uniqueIndex"`
	LastProcessedID int
}

type GotifyMsg struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	APPID   int    `json:"appid"`
}
