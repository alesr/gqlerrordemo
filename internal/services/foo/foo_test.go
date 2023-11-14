package foo

import (
	"testing"

	"github.com/alesr/gqlerrordemo/internal/repository/fakedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Repository Mock

var _ repo = (*repoMock)(nil)

type repoMock struct {
	CreateFunc func(id string) error
	GetFunc    func(id string) (string, error)
}

func (m *repoMock) Create(id string) error {
	return m.CreateFunc(id)
}

func (m *repoMock) Get(id string) (string, error) {
	return m.GetFunc(id)
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("should return foo", func(t *testing.T) {
		t.Parallel()

		repo := &repoMock{
			GetFunc: func(id string) (string, error) {
				return "foo", nil
			},
		}

		service := NewDefaultService(repo)

		observed, err := service.Get("dummy-id")
		require.NoError(t, err)

		require.Equal(t, "foo", observed.ID)
	})

	t.Run("should return foo not found error", func(t *testing.T) {
		t.Parallel()

		repo := &repoMock{
			GetFunc: func(id string) (string, error) {
				return "", fakedb.ErrNotFound
			},
		}

		service := NewDefaultService(repo)

		observed, err := service.Get("dummy-id")
		require.Error(t, err)

		require.Nil(t, observed)

		assert.ErrorIs(t, err, ErrFooNotFound)
	})
}
