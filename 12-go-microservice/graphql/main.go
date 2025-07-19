package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(
		cfg.AccountURL,
		cfg.CatalogURL,
		cfg.OrderURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new server with explicit transport configuration
	srv := handler.New(s.ToExecutableSchema())

	// Add the transports we want to support
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("Khoa Le", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
