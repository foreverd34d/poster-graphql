package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"

	"github.com/foreverd34d/poster-graphql/graph/model"
)

// Comments is the resolver for the comments field.
func (r *commentResolver) Comments(ctx context.Context, obj *model.Comment, offset *int, limit *int) ([]*model.Comment, error) {
	if len(obj.Comments) == 0 {
		return obj.Comments, nil
	}
	off := 0
	if offset != nil && *offset < len(obj.Comments) {
		off = *offset
	}
	lim := len(obj.Comments)
	if limit != nil && *limit <= len(obj.Comments) {
		lim = *limit
	}
	return obj.Comments[off:lim], nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	return r.service.CreatePost(ctx, input)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	comment, err := r.service.CreateComment(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, observer := range r.postObservers {
		if comment.PostID != nil && comment.PostID.String() == observer.PostID {
			observer.Ch <- comment
		}
	}

	return comment, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *model.Post, offset *int, limit *int) ([]*model.Comment, error) {
	if len(obj.Comments) == 0 {
		return obj.Comments, nil
	}
	off := 0
	if offset != nil && *offset < len(obj.Comments) {
		off = *offset
	}
	lim := len(obj.Comments)
	if limit != nil && *limit <= len(obj.Comments) {
		lim = *limit
	}
	return obj.Comments[off:lim], nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, offset *int, limit *int) ([]*model.Post, error) {
	return r.service.GetAllPosts(ctx, offset, limit)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	return r.service.GetPostByID(ctx, id)
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	newObserver := postObserver{
		Ch: make(chan *model.Comment),
		PostID: postID,
	}

	r.mu.Lock()
	r.postObservers = append(r.postObservers, newObserver)
	r.mu.Unlock()

	id := len(r.postObservers) - 1
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		r.postObservers = append(r.postObservers[:id], r.postObservers[:id+1]...)
		r.mu.Unlock()
	}()

	return r.postObservers[id].Ch, nil
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
