package inmem

import (
	"context"
	"errors"
	"slices"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/google/uuid"
)

var (
	ErrNotCommentable = errors.New("the post is not commentable")
	ErrNotFound = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type Repo struct {
	posts []*model.Post
}

func NewRepo() *Repo {
	return &Repo{make([]*model.Post, 0)}
}

func (r *Repo) CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error) {
	post := &model.Post{
		ID: uuid.New(),
		Title: newPost.Title,
		Author: newPost.Author,
		Content: newPost.Content,
		Commentable: newPost.Commentable,
	}
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *Repo) GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error) {
	off := 0
	if offset != nil && *offset < len(r.posts) {
		off = *offset
	}
	lim := len(r.posts)
	if limit != nil && *limit <= len(r.posts) {
		lim = *limit
	}

	if len(r.posts) == 0 {
		return r.posts, nil
	}
	return r.posts[off:lim], nil
}

func (r *Repo) GetPostById(ctx context.Context, id string) (*model.Post, error) {
	idx := slices.IndexFunc(r.posts, func(post *model.Post) bool {
		return post.ID.String() == id
	})
	if idx == -1 {
		return nil, ErrNotFound
	}
	return r.posts[idx], nil
}

func (r *Repo) CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error) {
	comm := &model.Comment{
		ID: uuid.New(),
		Author: newComment.Author,
		Content: newComment.Content,
	}

	if newComment.PostID != nil {
		idx := slices.IndexFunc(r.posts, func(post *model.Post) bool {
			return post.ID.String() == *newComment.PostID
		})
		if idx == -1 {
			return nil, ErrNotFound
		}
		if !r.posts[idx].Commentable {
			return nil, ErrNotCommentable
		}
		r.posts[idx].Comments = append(r.posts[idx].Comments, comm)
	} else if newComment.CommentID != nil {
		for i := range r.posts {
			r.posts[i].Comments = insertReply(r.posts[i].Comments, comm, *newComment.CommentID)
		}
	} else {
		return nil, ErrBadRequest
	}

	return comm, nil
}

func insertReply(comms []*model.Comment, newComm *model.Comment, commId string) []*model.Comment {
	if comms == nil {
		return nil
	}

	for i := range comms {
		if comms[i].ID.String() == commId {
			comms[i].Comments = append(comms[i].Comments, newComm)
		} else {
			comms[i].Comments = insertReply(comms[i].Comments, newComm, commId)
		}
	}

	return comms
}
