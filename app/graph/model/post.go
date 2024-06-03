package model

import "github.com/google/uuid"

type Post struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Content     string     `json:"content"`
	Commentable bool       `json:"commentable"`
	Comments    []*Comment `json:"comments,omitempty"`
}
