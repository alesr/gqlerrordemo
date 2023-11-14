package foo

import (
	"errors"
	"fmt"

	"github.com/alesr/gqlerrordemo/internal/repository/fakedb"
)

type repo interface {
	Create(id string) error
	Get(id string) (string, error)
}

// Foo is represents a foo domain object.
type Foo struct {
	ID string `json:"id"`
}

// DefaultService is the default implementation of the foo service.
type DefaultService struct {
	repo repo
}

// NewDefaultService creates a new instance of the default foo service.
func NewDefaultService(repo repo) *DefaultService {
	return &DefaultService{
		repo: repo,
	}
}

// Create creates a new foo with the given id.
func (s *DefaultService) Create(id string) (*Foo, error) {
	if id == "" {
		return nil, ErrFooMissingID
	}

	if err := s.repo.Create(id); err != nil {
		// In this case, we want to return a more user friendly error
		// since it is a usefull information for the user.
		if errors.Is(err, fakedb.ErrDuplicateEntry) {
			return nil, ErrFooAlreadyExists
		}
		return nil, fmt.Errorf("could not create foo: %w", err)
	}

	return &Foo{
		ID: id,
	}, nil
}

// Get gets a foo with the given id.
func (s *DefaultService) Get(id string) (*Foo, error) {
	if id == "" {
		return nil, ErrFooMissingID
	}

	name, err := s.repo.Get(id)
	if err != nil {
		// In this case, we want to return a more user friendly error
		// since it is a usefull information for the user.
		if errors.Is(err, fakedb.ErrNotFound) {
			return nil, ErrFooNotFound
		}
		return nil, fmt.Errorf("could not get foo: %w", err)
	}

	return &Foo{
		ID: name,
	}, nil
}
