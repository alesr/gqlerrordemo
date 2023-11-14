package fakedb

import "errors"

var (
	// Enumerate all possible errors that can be returned by the repository.

	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrNotFound       = errors.New("not found")
)

// FakeDB is a fake implementation of a database.
type FakeDB struct{}

// New creates a new instance of the fake database.
func New() *FakeDB {
	return &FakeDB{}
}

// Create creates a new foo with the given id.
func (db *FakeDB) Create(id string) error {
	// Simulate an already exists error.
	if id == "alreadyexists" {
		return ErrDuplicateEntry
	}
	return nil
}

// Get gets a foo with the given id.
func (db *FakeDB) Get(id string) (string, error) {
	// Simulate a not found error.
	if id == "notfound" {
		return "", ErrNotFound
	}
	return "foo", nil
}
