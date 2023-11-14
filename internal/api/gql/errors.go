package gql

import (
	"context"
	"errors"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/alesr/gqlerrordemo/internal/services/foo"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	// Enumerate all possible errors that can be returned by the API.

	errFooAlreadyExist = gqlerror.Error{
		Message: "Foo already exists",
		Extensions: map[string]any{
			"code": "FOO_ALREADY_EXISTS",
		},
	}

	errFooMissingID = gqlerror.Error{
		Message: "Foo missing id",
		Extensions: map[string]any{
			"code": "FOO_MISSING_ID",
		},
	}

	errFooNotFound = gqlerror.Error{
		Message: "Foo not found",
		Extensions: map[string]any{
			"code": "FOO_NOT_FOUND",
		},
	}

	// Default error to return if no matching error is found.
	errInternalServer = gqlerror.Error{
		Message: "Internal server error",
		Extensions: map[string]any{
			"code": "INTERNAL_SERVER_ERROR",
		},
	}
)

var ErrorPresenter = func(logger *slog.Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		if err == nil {
			return &errInternalServer
		}

		logger.Info("api error", slog.String("original_error", err.Error()))

		apiErr := parseError(ctx, err)

		// Add the path to the error.
		apiErr.Path = graphql.GetPath(ctx)

		return &apiErr
	}
}

// parseError returns the corresponding API error for the given error.
// If no matching error is found, it returns the default error.
func parseError(ctx context.Context, err error) gqlerror.Error {
	switch {
	case errors.Is(err, foo.ErrFooAlreadyExists):
		return errFooAlreadyExist
	case errors.Is(err, foo.ErrFooMissingID):
		return errFooMissingID
	case errors.Is(err, foo.ErrFooNotFound):
		return errFooNotFound
	default:
		return errInternalServer
	}
}
