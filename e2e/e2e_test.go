package e2e

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alesr/gqlerrordemo/internal/api/gql"
	"github.com/alesr/gqlerrordemo/internal/repository/fakedb"
	"github.com/alesr/gqlerrordemo/internal/services/foo"
	"github.com/stretchr/testify/require"
)

func TestGetFoo(t *testing.T) {
	srv := setupHelper(t)
	defer teardownHelper(t, srv)

	t.Run("GetFoo success", func(t *testing.T) {
		query := `query {
			getFoo(id: "dummy-id") {
				id
			}
		}`

		expected := `{"data":{"getFoo":{"id":"foo"}}}`

		actual := doRequestHelper(t, query)

		require.JSONEq(t, expected, actual)
	})

	t.Run("GetFoo error", func(t *testing.T) {
		query := `query {
			getFoo(id: "notfound") {
				id
			}
		}`

		expected := `{
			"errors": [
				{
					"message": "Foo not found",
					"path": ["getFoo"],
					"extensions": {
						"code": "FOO_NOT_FOUND"
					}
				}
			],
			"data": {
				"getFoo": null
			}
		}`

		actual := doRequestHelper(t, query)

		require.JSONEq(t, expected, actual)
	})

	t.Run("CreateFoo success", func(t *testing.T) {
		query := `mutation {
			createFoo(id: "dummy-id") {
				id
			}
		}`

		expected := `{"data":{"createFoo":{"id":"dummy-id"}}}`

		actual := doRequestHelper(t, query)

		require.JSONEq(t, expected, actual)
	})

	t.Run("CreateFoo error", func(t *testing.T) {
		query := `mutation {
			createFoo(id: "alreadyexists") {
				id
			}
		}`

		expected := `{
			"errors": [
				{
					"message": "Foo already exists",
					"path": ["createFoo"],
					"extensions": {
						"code": "FOO_ALREADY_EXISTS"
					}
				}
			],
			"data": {
				"createFoo": null
			}
		}`

		actual := doRequestHelper(t, query)

		require.JSONEq(t, expected, actual)
	})
}

func setupHelper(t *testing.T) *http.Server {
	t.Helper()

	port := "8080"

	noopLogger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	fooRepo := fakedb.New()

	fooService := foo.NewDefaultService(fooRepo)

	resolvers := gql.NewResolvers(noopLogger, fooService)

	schema := gql.NewExecutableSchema(
		gql.Config{
			Resolvers: resolvers,
		},
	)

	handler := handler.NewDefaultServer(schema)

	handler.SetErrorPresenter(gql.ErrorPresenter(noopLogger))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handler)

	srv := &http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}

func teardownHelper(t *testing.T, srv *http.Server) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, srv.Shutdown(ctx))
}

func doRequestHelper(t *testing.T, gqlOperation string) string {
	t.Helper()

	payload := strings.NewReader(fmt.Sprintf(`{"query": %q}`, gqlOperation))

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/query", payload)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return string(body)
}
