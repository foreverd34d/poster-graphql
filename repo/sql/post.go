package sql

import (
	"context"
	"database/sql"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/jmoiron/sqlx"
)

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db}
}

func (r *PostRepo) CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error) {
	post := new(model.Post)
	query := "INSERT INTO posts (title, author, content, commentable) VALUES ($1, $2, $3, $4) RETURNING id, title, author, content, commentable"
	if err := r.db.GetContext(ctx, post, query, newPost.Title, newPost.Author, newPost.Content, newPost.Commentable); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepo) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post
	query := "SELECT id, title, author, content, commentable FROM posts"
	if err := r.db.SelectContext(ctx, &posts, query); err != nil {
		return nil, err
	}

	var err error
	commsQuery := `SELECT id, author, content, post_id FROM comments WHERE post_id = $1`
	for i := range posts {
		posts[i].Comments, err = retrieveReplies(ctx, r.db, commsQuery, posts[i].ID)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetPostById(ctx context.Context, id string) (*model.Post, error) {
	post := new(model.Post)
	postQuery := "SELECT * FROM posts WHERE id = $1"
	if err := r.db.GetContext(ctx, post, postQuery, id); err != nil {
		return nil, err
	}

	var err error
	commsQuery := `SELECT id, author, content, post_id FROM comments WHERE post_id = $1`
	post.Comments, err = retrieveReplies(ctx, r.db, commsQuery, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func retrieveReplies(ctx context.Context, db *sqlx.DB, query string, args ...any) ([]*model.Comment, error) {
	var replies []*model.Comment
	if err := db.SelectContext(ctx, &replies, query, args...); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	for i := range replies {
		var err error
		replQuery := `SELECT id, author, content, post_id FROM comments WHERE parent_comment_id = $1`
		replies[i].Comments, err = retrieveReplies(ctx, db, replQuery, replies[i].ID.String())
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}

	return replies, nil
}

// func retrieveRepliesAsync(ctx context.Context, db *sqlx.DB, parentCommentId string, ch chan *model.Comment) {
// 	var replies []*model.Comment
// 	query := "SELECT * FROM comment WHERE parent_comment_id = $1"
// 	if err := db.SelectContext(ctx, &replies, query, parentCommentId); err != nil && err != sql.ErrNoRows {
// 		return
// 	}
//
// 	replCh := make(chan *model.Comment, len(replies))
// 	for i := range replies {
// 		go retrieveRepliesAsync(ctx, db, replies[i].ID.String(), replCh)
// 	}
// }
