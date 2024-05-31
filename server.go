package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/foreverd34d/poster-graphql/graph"
	"github.com/foreverd34d/poster-graphql/repo"
	"github.com/foreverd34d/poster-graphql/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var r *repo.Repo
	dbConnection := os.Getenv("DBCONNECT")
	if dbConnection == "" {
		r = repo.NewInMemRepo()
	} else {
		db, err  := sqlx.Connect("postgres", dbConnection)
		if err != nil {
			log.Fatalln(err)
		}
		r = repo.NewSqlRepo(db)
	}

	s := service.NewService(r)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(s)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
