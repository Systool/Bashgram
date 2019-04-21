package main

import "encoding/json"

type APIResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from"`
	Date      int    `json:"date"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName     string `json:"username"`
}
