package foo

import "fmt"

var (
	// Enumerate all possible errors returned by the service.
	// Exported errors are considered part of the service API,
	// and are used in the transport layer to return the correct
	// error to the client.
	//
	// Errors that are not expected to be visible to the client
	// should not be exported.

	ErrFooAlreadyExists = fmt.Errorf("foo already exists")
	ErrFooMissingID     = fmt.Errorf("foo missing id")
	ErrFooNotFound      = fmt.Errorf("foo not found")
)
