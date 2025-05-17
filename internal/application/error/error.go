package errors

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	appcontext "GabrielChaves1/notilify/internal/application/context"
)

const (
	ValidationError      string = "VALIDATION_ERROR"
	NotFoundError        string = "NOT_FOUND_ERROR"
	AlreadyExistsError   string = "ALREADY_EXISTS_ERROR"
	InternalServerError  string = "INTERNAL_SERVER_ERROR"
	UnauthorizedError    string = "UNAUTHORIZED_ERROR"
	ForbiddenError       string = "FORBIDDEN_ERROR"
	BadRequestError      string = "BAD_REQUEST_ERROR"
	ConflictError        string = "CONFLICT_ERROR"
	ExternalServiceError string = "EXTERNAL_SERVICE_ERROR"
)

type AppError struct {
	Code      string
	Message   string
	Context   map[string]interface{}
	Timestamp time.Time
	Cause     error
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) Is(target error) bool {
	if targetErr, ok := target.(*AppError); ok {
		return e.Code == targetErr.Code
	}
	return false
}

func NewValidationError(ctx context.Context, validationErrors *ValidationErrors) *AppError {
	currentTime, err := appcontext.ExtractCurrentTimeFromContext(ctx)
	if err != nil {
		currentTime = time.Time{}
	}

	return &AppError{
		Code:      ValidationError,
		Message:   "One or more validation errors occurred",
		Context:   map[string]interface{}{},
		Timestamp: currentTime,
		Cause:     validationErrors,
	}
}

func NewInternalServerError(ctx context.Context, reason string) *AppError {
	currentTime, err := appcontext.ExtractCurrentTimeFromContext(ctx)
	if err != nil {
		currentTime = time.Time{}
	}

	return &AppError{
		Code:    InternalServerError,
		Message: "Internal server error",
		Context: map[string]interface{}{
			"reason": reason,
		},
		Timestamp: currentTime,
	}
}

func NewInvalidContextError(ctx context.Context, fieldName string, reason string) *AppError {
	currentTime, err := appcontext.ExtractCurrentTimeFromContext(ctx)
	if err != nil {
		currentTime = time.Time{}
	}

	return &AppError{
		Code:    InternalServerError,
		Message: fmt.Sprintf("Valor de '%s' inválido: %s", fieldName, reason),
		Context: map[string]interface{}{
			"fieldName": fieldName,
			"reason":    reason,
		},
		Timestamp: currentTime,
	}
}

func IsNotFound(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == NotFoundError
}

func IsValidationError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == ValidationError
}

func IsConflict(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == ConflictError
}

func IsSystemError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && (appErr.Code == InternalServerError || strings.HasPrefix(appErr.Code, "SYS_"))
}

func IsAuthorizationError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == UnauthorizedError
}

func IsExternalServiceError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == ExternalServiceError
}

func GetValidationErrors(err error) *ValidationErrors {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Code == ValidationError {
		var validationErrs *ValidationErrors
		if errors.As(appErr.Cause, &validationErrs) {
			return validationErrs
		}
	}
	return nil
}

func GetErrorContext(err error) (map[string]interface{}, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Context, true
	}
	return nil, false
}

