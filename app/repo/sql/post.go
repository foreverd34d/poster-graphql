package sql

import (
	"context"
	"database/sql"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db}
}

func (r *postRepo) CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error) {
	post := new(model.Post)
	query := "INSERT INTO posts (title, author, content, commentable) VALUES ($1, $2, $3, $4) RETURNING id, title, author, content, commentable"
	if err := r.db.GetContext(ctx, post, query, newPost.Title, newPost.Author, newPost.Content, newPost.Commentable); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepo) GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error) {
	var posts []*model.Post
	query := "SELECT id, title, author, content, commentable FROM posts OFFSET $1 LIMIT $2"
	if err := r.db.SelectContext(ctx, &posts, query, *offset, *limit); err != nil {
		return nil, err
	}

	commsQuery := `SELECT id, author, content, post_id FROM comments WHERE post_id = $1`
	ch := make(chan chComment, len(posts))
	for i := range posts {
		go retrieveReplies(ctx, r.db, i, ch, commsQuery, posts[i].ID)
	}
	for range posts {
		comment := <-ch
		posts[comment.idx].Comments = comment.comms
	}

	return posts, nil
}

func (r *postRepo) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	post := new(model.Post)
	postQuery := "SELECT * FROM posts WHERE id = $1"
	if err := r.db.GetContext(ctx, post, postQuery, id); err != nil {
		return nil, err
	}

	commsQuery := `SELECT id, author, content, post_id FROM comments WHERE post_id = $1`
	ch := make(chan chComment)
	retrieveReplies(ctx, r.db, 0, ch, commsQuery, id)
	comments := <- ch
	post.Comments = comments.comms

	return post, nil
}

type chComment struct {
	idx  int
	comms []*model.Comment
}

func retrieveReplies(ctx context.Context, db *sqlx.DB, replyIdx int, ch chan chComment, query string, args ...any) {
	var replies []*model.Comment
	if err := db.SelectContext(ctx, &replies, query, args...); err != nil && err != sql.ErrNoRows {
		ch <- chComment{replyIdx, nil}
		return
	}

	if len(replies) == 0 {
		ch <- chComment{replyIdx, nil}
		return
	}

	replCh := make(chan chComment, len(replies))
	for i := range replies {
		replQuery := "SELECT * FROM comment WHERE parent_comment_id = $1"
		go retrieveReplies(ctx, db, i, replCh, replQuery, replies[i].ID.String())
	}

	for range replies {
		reply := <-replCh
		replies[reply.idx].Comments = reply.comms
	}

	ch <- chComment{replyIdx, replies}
}
