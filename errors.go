package vipertemplate

import (
	"errors"
	"fmt"
)

func newError(key string, err error) error {
	return fmt.Errorf("non parsable template for the key '%s': %w", key, err)
}

var (
	// ErrCircularDependency is an error returned if the circular dependency detected
	ErrCircularDependency = errors.New("circular dependency")
	// ErrNonStringValue is an error returned if the requested value is not a string
	ErrNonStringValue = errors.New("non string value")
)
