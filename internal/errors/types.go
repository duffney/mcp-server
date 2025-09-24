package errors

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

type ErrorCategory string

const (
	ValidationError     ErrorCategory = "validation"
	ExecutionError      ErrorCategory = "execution"
	SystemError         ErrorCategory = "system"
	NetworkError        ErrorCategory = "network"
	AuthenticationError ErrorCategory = "authentication"
)

type CopaceticError struct {
	Category ErrorCategory
	Op       string // Operation that failed
	Image    string // Container image (when applicable)
	Cause    error  // Underlying error
}

func (e *CopaceticError) Error() string {
	if e.Image != "" {
		return fmt.Sprintf("%s failed for image %s: %v", e.Op, e.Image, e.Cause)
	}
	return fmt.Sprintf("%s failed: %v", e.Op, e.Cause)
}

func (e *CopaceticError) Unwrap() error { return e.Cause }

func NewValidationError(op, image string, cause error) *CopaceticError {
	return &CopaceticError{
		Category: ValidationError,
		Op:       op,
		Image:    image,
		Cause:    errors.Wrap(cause, op),
	}
}

func NewExecutionError(op, image string, cause error) *CopaceticError {
	return &CopaceticError{
		Category: ExecutionError,
		Op:       op,
		Image:    image,
		Cause:    errors.Wrap(cause, op),
	}
}

func NewSystemError(op, image string, cause error) *CopaceticError {
	return &CopaceticError{
		Category: SystemError,
		Op:       op,
		Image:    image,
		Cause:    errors.Wrap(cause, op),
	}
}

func NewAuthenticationError(op, image string, cause error) *CopaceticError {
	return &CopaceticError{
		Category: AuthenticationError,
		Op:       op,
		Image:    image,
		Cause:    errors.Wrap(cause, op),
	}
}
