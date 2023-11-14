package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alesr/gqlerrordemo/internal/api/gql"
	graph "github.com/alesr/gqlerrordemo/internal/api/gql"
	"github.com/alesr/gqlerrordemo/internal/repository/fakedb"
	"github.com/alesr/gqlerrordemo/internal/services/foo"
)

const defaultPort = "8080"

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				AddSource: true,
			},
		),
	)

	// Init repository
	fooRepo := fakedb.New()

	// Init service
	fooService := foo.NewDefaultService(fooRepo)

	// Init resolvers
	resolvers := graph.NewResolvers(logger, fooService)

	// Init graphql schema
	schema := gql.NewExecutableSchema(
		gql.Config{
			Resolvers: resolvers,
		},
	)

	// Init graphql server
	gqlHandler := handler.NewDefaultServer(schema)

	// Set error presenter
	gqlHandler.SetErrorPresenter(gql.ErrorPresenter(logger))

	// Init http server
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", gqlHandler)

	srv := &http.Server{
		Addr: ":" + defaultPort,
	}

	// Run http server and wait for cancel signal

	logger.Info(
		"starting gqlerrordemo server",
		slog.String("GraphQL playground", fmt.Sprintf("http://localhost:%s/", defaultPort)),
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Error(err.Error())
			}
			logger.Info("gqlerrordemo server closed")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	// Shutdown http server

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("gqlerrordemo server stopped")
}
