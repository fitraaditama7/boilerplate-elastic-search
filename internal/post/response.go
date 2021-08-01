package post

import "time"

type createResponse struct {
	ID string `json:"id"`
}

type findResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
