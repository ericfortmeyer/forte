package archive

import "errors"

// SkippableError indicates an error that doesn't require deployment to fail.
// Used for validation failures or missing archives—extraction can be skipped.
type SkippableError struct {
	message string
}

func (e *SkippableError) Error() string {
	return e.message
}

// NewSkippableError creates a new SkippableError.
func NewSkippableError(msg string) *SkippableError {
	return &SkippableError{message: msg}
}

// IsSkippable checks if an error should not fail deployment.
func IsSkippable(err error) bool {
	if err == nil {
		return false
	}
	var se *SkippableError
	return errors.As(err, &se)
}
