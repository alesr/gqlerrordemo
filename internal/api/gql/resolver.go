package gql

import (
	"log/slog"

	"github.com/alesr/gqlerrordemo/internal/services/foo"
)

type fooService interface {
	Create(id string) (*foo.Foo, error)
	Get(id string) (*foo.Foo, error)
}

// Resolver is the root resolver for the GraphQL API
// it contains all the services that will be used by the resolvers.
type Resolver struct {
	logger     *slog.Logger
	FooService fooService
}

// NewResolvers returns a new instance of the root resolver.
func NewResolvers(logger *slog.Logger, fooService fooService) *Resolver {
	return &Resolver{
		logger:     logger,
		FooService: fooService,
	}
}
