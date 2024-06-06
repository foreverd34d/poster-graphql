package sql

import (
	"context"
	"database/sql"
	"sync"

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

	commsQuery := "SELECT id, author, content, post_id FROM comments WHERE post_id = $1"
	var wg sync.WaitGroup
	wg.Add(len(posts))
	for _, post := range posts {
		go func(p *model.Post) {
			defer wg.Done()
			for comment := range retrieveComments(ctx, r.db, commsQuery, p.ID) {
				p.Comments = append(p.Comments, comment)
			}
		}(post)
	}
	wg.Wait()

	return posts, nil
}

func (r *postRepo) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	post := new(model.Post)
	postQuery := "SELECT * FROM posts WHERE id = $1"
	if err := r.db.GetContext(ctx, post, postQuery, id); err != nil {
		return nil, err
	}

	commsQuery := "SELECT id, author, content, post_id FROM comments WHERE post_id = $1"
	for comment := range retrieveComments(ctx, r.db, commsQuery, id) {
		post.Comments = append(post.Comments, comment)
	}

	return post, nil
}

func retrieveComments(ctx context.Context, db *sqlx.DB, query string, args ...any) <-chan *model.Comment {
	ch := make(chan *model.Comment)
	go func() {
		defer close(ch)
		var comments []*model.Comment
		if err := db.SelectContext(ctx, &comments, query, args...); err != nil && err != sql.ErrNoRows {
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(comments))
		for _, comment := range comments {
			repliesQuery := "SELECT id, author, content, post_id FROM comments WHERE parent_comment_id = $1"
			go func(c *model.Comment) {
				defer wg.Done()
				for reply := range retrieveComments(ctx, db, repliesQuery, c.ID) {
					c.Comments = append(c.Comments, reply)
				}
				ch <- c
			}(comment)
		}
		wg.Wait()
	}()
	return ch
}
