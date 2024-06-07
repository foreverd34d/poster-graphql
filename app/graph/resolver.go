package graph

import (
	"sync"

	"github.com/foreverd34d/poster-graphql/graph/model"
	"github.com/foreverd34d/poster-graphql/service"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type postObserver struct {
	Ch chan *model.Comment
	PostID string
}

type Resolver struct {
	service *service.Service
	postObservers []postObserver
	mu sync.Mutex
}

func NewResolver(s *service.Service) *Resolver {
	return &Resolver{
		service: s,
		postObservers: make([]postObserver, 0),
	}
}
