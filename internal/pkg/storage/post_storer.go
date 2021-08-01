package storage

import (
	"context"
	"time"
)

type PostStorer interface {
	Insert(ctx context.Context, post Post) error
	Update(ctx context.Context, post Post) error
	Delete(ctx context.Context, id string) error
	FindOne(ctx context.Context, id string) (Post, error)
	Search(ctx context.Context, keyword string) (interface{}, error)
}

type Post struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Text      string     `json:"text"`
	Tags      []string   `json:"tags"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
