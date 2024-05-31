package service

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo"
)

type PostService struct {
	repo repo.Post
}

func NewPostService(repo repo.Post) *PostService {
	return &PostService{repo}
}

func (s *PostService) CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error) {
	return s.repo.CreatePost(ctx, newPost)
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	return s.repo.GetAllPosts(ctx)
}

func (s *PostService) GetPostById(ctx context.Context, id string) (*model.Post, error) {
	return s.repo.GetPostById(ctx, id)
}
