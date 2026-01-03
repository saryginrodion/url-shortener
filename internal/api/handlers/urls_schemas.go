package handlers

import "time"

type UrlCreateRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url" validate:"required,url"`
}

type UrlDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}
