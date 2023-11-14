package gql

import (
	"context"
	"testing"

	"github.com/alesr/gqlerrordemo/internal/services/foo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Foo Service Mock

var _ fooService = (*fooServiceMock)(nil)

type fooServiceMock struct {
	CreateFunc func(id string) (*foo.Foo, error)
	GetFunc    func(id string) (*foo.Foo, error)
}

func (m *fooServiceMock) Create(id string) (*foo.Foo, error) {
	return m.CreateFunc(id)
}

func (m *fooServiceMock) Get(id string) (*foo.Foo, error) {
	return m.GetFunc(id)
}

func TestCreateFooResolver(t *testing.T) {
	t.Parallel()

	t.Run("should create a foo with the given id", func(t *testing.T) {
		t.Parallel()

		givenID := "dummy-id"

		fooService := &fooServiceMock{
			CreateFunc: func(id string) (*foo.Foo, error) {
				return &foo.Foo{
					ID: givenID,
				}, nil
			},
		}

		resolver := &Resolver{
			FooService: fooService,
		}

		result, err := resolver.Mutation().CreateFoo(context.TODO(), givenID)
		require.NoError(t, err)

		require.Equal(t, givenID, result.ID)
	})

	t.Run("should return a generic error if the service return an unexpected error", func(t *testing.T) {
		t.Parallel()

		givenServiceError := assert.AnError

		fooService := &fooServiceMock{
			CreateFunc: func(id string) (*foo.Foo, error) {
				return nil, givenServiceError
			},
		}

		resolver := &Resolver{
			FooService: fooService,
		}

		result, err := resolver.Mutation().CreateFoo(context.TODO(), "dummy-id")
		require.Error(t, err)

		require.Nil(t, result)
		require.Equal(t, givenServiceError, err)
	})
}
