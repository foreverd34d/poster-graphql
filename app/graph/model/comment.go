package model

import "github.com/google/uuid"

type Comment struct {
	ID       uuid.UUID  `json:"id"`
	Author   string     `json:"author"`
	Content  string     `json:"content"`
	PostID   *uuid.UUID `json:"postId,omitempty" db:"post_id"`
	Comments []*Comment `json:"comments,omitempty"`
}
