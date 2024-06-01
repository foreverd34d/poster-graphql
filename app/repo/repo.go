package repo

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo/inmem"
	"github.com/foreverd34d/poster-graphql/repo/sql"
	"github.com/jmoiron/sqlx"
)

type Post interface {
	CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error)
	GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error)
	GetPostById(ctx context.Context, id string) (*model.Post, error)
}

type Comment interface {
	CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error)
}

type Repo struct {
	Post
	Comment
}

func NewSqlRepo(db *sqlx.DB) *Repo {
	return &Repo{
		Post: sql.NewPostRepo(db),
		Comment: sql.NewCommentRepo(db),
	}
}

func NewInMemRepo() *Repo {
	repo := inmem.NewRepo()
	return &Repo{
		Post: repo,
		Comment: repo,
	}
}
