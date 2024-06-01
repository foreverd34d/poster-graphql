package service

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo"
)

type Post interface {
	CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error)
	GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error)
	GetPostById(ctx context.Context, id string) (*model.Post, error)
}

type Comment interface {
	CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error)
}

type Service struct {
	Post
	Comment
}

func NewService(repo *repo.Repo) *Service {
	return &Service{
		Post: NewPostService(repo),
		Comment: NewCommentService(repo),
	}
}
