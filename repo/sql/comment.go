package sql

import (
	"context"
	"errors"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/jmoiron/sqlx"
)

type CommentRepo struct {
	db *sqlx.DB
}

var (
	ErrNotCommentable = errors.New("the post is not commentable")
)

func NewCommentRepo(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{db}
}

func (r *CommentRepo) CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error) {
	if newComment.PostID != nil {
		var commentable bool
		query := "SELECT commentable FROM posts WHERE id = $1"
		if err := r.db.GetContext(ctx, &commentable, query, *newComment.PostID); err != nil {
			return nil, err
		}
		if !commentable {
			return nil, ErrNotCommentable
		}
	}

	comment := new(model.Comment)
	query := "INSERT INTO comments (author, content, parent_comment_id, post_id) VALUES ($1, $2, $3, $4) RETURNING id, author, content, post_id"
	if err := r.db.GetContext(ctx, comment, query, newComment.Author, newComment.Content, newComment.CommentID, newComment.PostID); err != nil {
		return nil, err
	}
	return comment, nil
}
