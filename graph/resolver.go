package graph

import "github.com/foreverd34d/poster-graphql/service"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	service *service.Service
}

func NewResolver(s *service.Service) *Resolver {
	return &Resolver{s}
}
