package models

type Link struct {
	ID          int
	ActiveLink  string `json:"active_link"`
	HistoryLink string `json:"history_link"`
}
