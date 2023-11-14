package gql

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/alesr/gqlerrordemo/internal/services/foo"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var noopLogger = slog.New(slog.NewJSONHandler(io.Discard, nil))

func gqlCtxHelper(t *testing.T, ctx context.Context, path string) context.Context {
	t.Helper()
	return graphql.WithPathContext(ctx, graphql.NewPathWithField(path))
}

func TestErrorPresenter(t *testing.T) {
	t.Parallel()

	t.Run("should return defined api error if matching with the given error", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name       string
			givenError error
			expected   *gqlerror.Error
		}{
			{
				name:       "nil",
				givenError: nil,
				expected:   &errInternalServer,
			},
			{
				name:       "should return internal server error as it is not an api error",
				givenError: errors.New("foo"),
				expected:   &errInternalServer,
			},
			{
				name:       "should return foo already exist error",
				givenError: foo.ErrFooAlreadyExists,
				expected:   &errFooAlreadyExist,
			},
			{
				name:       "should return foo missing id error",
				givenError: foo.ErrFooMissingID,
				expected:   &errFooMissingID,
			},
			{
				name:       "should return foo not found error",
				givenError: foo.ErrFooNotFound,
				expected:   &errFooNotFound,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				observed := ErrorPresenter(noopLogger)(context.Background(), tc.givenError)

				assert.Equal(t, tc.expected, observed)
			})
		}
	})

	t.Run("should return the error with the path from the context", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name         string
			givenContext context.Context
			expectedPath string
		}{
			{
				name:         "no path",
				givenContext: context.Background(),
				expectedPath: "",
			},
			{
				name:         "path with field",
				givenContext: gqlCtxHelper(t, context.Background(), "foo"),
				expectedPath: "foo",
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				observed := ErrorPresenter(noopLogger)(tc.givenContext, foo.ErrFooAlreadyExists)

				assert.Equal(t, tc.expectedPath, observed.Path.String())
			})
		}
	})
}
