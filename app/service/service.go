package service

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo"
)

// Post defines methods for creating and getting posts with all comments.
type Post interface {
	// CreatePost creates new post and returns it with ID and empty comments.
	CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error)

	// GetAllPosts gets all posts with all comments within [offset, limit) range.
	GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error)

	// GetPostByID gets the post with all comments by ID.
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
}

// Comment defines methods for creating comments to the post or another comment.
type Comment interface {
	// CreateComment creates new comment and returns it with ID/
	CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error)
}

// Service groups all methods for working with posts and comments.
type Service struct {
	Post
	Comment
}

// NewService creates an instance of service which utilize passed repository.
func NewService(repo *repo.Repo) *Service {
	return &Service{
		Post: NewPostService(repo),
		Comment: NewCommentService(repo),
	}
}
