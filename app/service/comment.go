package service

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo"
)

type CommentService struct {
	repo repo.Comment
}

func NewCommentService(repo repo.Comment) *CommentService {
	return &CommentService{repo}
}

func (s *CommentService) CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error) {
	return s.repo.CreateComment(ctx, newComment)
}
