package service

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo"
)

type postService struct {
	repo repo.Post
}

func NewPostService(repo repo.Post) *postService {
	return &postService{repo}
}

func (s *postService) CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error) {
	return s.repo.CreatePost(ctx, newPost)
}

func (s *postService) GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error) {
	return s.repo.GetAllPosts(ctx, offset, limit)
}

func (s *postService) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	return s.repo.GetPostByID(ctx, id)
}
