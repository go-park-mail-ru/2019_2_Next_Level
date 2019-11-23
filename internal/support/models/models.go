package models

import "time"

type Chat struct {
	Id string `json:"id"`
	StartDate time.Time`json:"startDate"`
	IsOpen bool `json:"isOpen"`
	Theme string `json:"theme"`
	Messages []Message `json:"messages,omitempty"`
	UserId string `json:"-"`
	SupportId []string `json:"-"`
}

type Message struct {
	Id string `json:"id"`
	UserNick string `json:"userNick"`
	Sent time.Time `json:"sent"`
	WasRead bool `json:"wasRead"`
	Body string `json:"body"`
	ChatId string `json:"-"`
}
