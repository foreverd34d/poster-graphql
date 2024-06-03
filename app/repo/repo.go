package repo

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/repo/inmem"
	"github.com/foreverd34d/poster-graphql/repo/sql"
	"github.com/jmoiron/sqlx"
)

// Post defines methods for storing and getting posts.
type Post interface {
	// CreatePost stores new post and returns it with ID and empty comments.
	CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error)

	// GetAllPosts gets all posts with all comments within [offset, limit) range.
	GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error)

	// GetPostByID gets the post with all comments by ID.
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
}

// Comment defines methods for storing comments.
type Comment interface {
	// CreateComment stores new comment and returns it with ID.
	// Comment's attachment to the post or other comment id defined by setting PostID or CommentID.
	// If post is not commentable, returns an error.
	CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error)
}

// Repo groups all methods for working with posts and comments
type Repo struct {
	Post
	Comment
}

// NewSqlRepo returns an instance of SQL-based database.
func NewSqlRepo(db *sqlx.DB) *Repo {
	return &Repo{
		Post:    sql.NewPostRepo(db),
		Comment: sql.NewCommentRepo(db),
	}
}

// NewInMemRepo returns an instance of in-memory based non-persistent database.
func NewInMemRepo() *Repo {
	repo := inmem.NewRepo()
	return &Repo{
		Post:    repo,
		Comment: repo,
	}
}
