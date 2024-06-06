package inmem

import (
	"context"
	"errors"
	"slices"
	"sync"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/google/uuid"
)

var (
	ErrNotCommentable = errors.New("the post is not commentable")
	ErrNotFound       = errors.New("not found")
	ErrBadRequest     = errors.New("bad request")
)

type inmemRepo struct {
	posts []*model.Post
}

func NewRepo() *inmemRepo {
	return &inmemRepo{make([]*model.Post, 0)}
}

func (r *inmemRepo) CreatePost(ctx context.Context, newPost model.NewPost) (*model.Post, error) {
	post := &model.Post{
		ID:          uuid.New(),
		Title:       newPost.Title,
		Author:      newPost.Author,
		Content:     newPost.Content,
		Commentable: newPost.Commentable,
	}
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *inmemRepo) GetAllPosts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error) {
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

func (r *inmemRepo) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	idx := slices.IndexFunc(r.posts, func(post *model.Post) bool {
		return post.ID.String() == id
	})
	if idx == -1 {
		return nil, ErrNotFound
	}
	return r.posts[idx], nil
}

func (r *inmemRepo) CreateComment(ctx context.Context, newComment model.NewComment) (*model.Comment, error) {
	comm := &model.Comment{
		ID:      uuid.New(),
		Author:  newComment.Author,
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
		var wg sync.WaitGroup
		wg.Add(len(r.posts))
		for _, post := range r.posts {
			go func(p *model.Post) {
				defer wg.Done()
				insertComment(p.Comments, comm, *newComment.CommentID)
			}(post)
		}
		wg.Wait()
	} else {
		return nil, ErrBadRequest
	}

	return comm, nil
}

func insertComment(comments []*model.Comment, newComment *model.Comment, parentId string) {
	if comments == nil {
		return
	}

	for _, comment := range comments {
		if comment.ID.String() == parentId {
			comment.Comments = append(comment.Comments, newComment)
			return
		} else {
			go insertComment(comment.Comments, newComment, parentId)
		}
	}
}
