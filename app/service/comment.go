package service

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo"
)

type commentService struct {
	repo repo.Comment
}

func NewCommentService(repo repo.Comment) *commentService {
	return &commentService{repo}
}

func (s *commentService) CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error) {
	return s.repo.CreateComment(ctx, newComment)
}
