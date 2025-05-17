package middleware

import (
	apperrors "GabrielChaves1/notilify/internal/application/error"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Detail  string                 `json:"detail,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	Status  int                    `json:"-"`
}

type APIErrorResponse struct {
	Errors []APIError             `json:"errors"`
	Meta   map[string]interface{} `json:"meta"`
}

var JSONErrors = []error{
	&json.SyntaxError{},
	&json.UnmarshalTypeError{},
	&json.InvalidUnmarshalError{},
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			requestIDStr := requestid.Get(c)

			apiErrors := mapToAPIErrors(err)

			status := http.StatusInternalServerError
			if len(apiErrors) > 0 {
				status = apiErrors[0].Status
			}

			response := APIErrorResponse{
				Errors: apiErrors,
				Meta: map[string]interface{}{
					"timestamp":  time.Now(),
					"request_id": requestIDStr,
				},
			}

			c.JSON(status, response)
			c.Abort()
		}
	}
}

func HandleError(c *gin.Context, err error) {
	_ = c.Error(err)

	apiErrs := mapToAPIErrors(err)
	status := http.StatusInternalServerError

	if len(apiErrs) > 0 {
		status = apiErrs[0].Status
	}

	response := APIErrorResponse{
		Errors: apiErrs,
		Meta: map[string]interface{}{
			"request_id": requestid.Get(c),
			"timestamp":  time.Now(),
		},
	}

	c.JSON(status, response)
	c.Abort()
}

func mapToAPIErrors(err error) []APIError {
	if apperrors.IsValidationError(err) {
		validationErrs := apperrors.GetValidationErrors(err)
		if validationErrs != nil && validationErrs.HasErrors() {
			return mapValidationErrorsToAPIErrors(validationErrs)
		}
	}

	return []APIError{mapErrorToAPIError(err)}
}

func mapValidationErrorsToAPIErrors(validationErrors *apperrors.ValidationErrors) []APIError {
	if validationErrors == nil || !validationErrors.HasErrors() {
		return []APIError{}
	}

	apiErrors := make([]APIError, 0, len(*validationErrors))

	for _, validationErr := range *validationErrors {
		apiErrors = append(apiErrors, APIError{
			Status:  http.StatusBadRequest,
			Code:    apperrors.ValidationError,
			Message: "Validation error",
			Detail:  formatValidationErrorDetail(validationErr),
			Context: createValidationErrorContext(validationErr),
		})
	}

	return apiErrors
}

func mapErrorToAPIError(err error) APIError {
	switch {
	case isJSONError(err):
		return APIError{
			Status:  http.StatusBadRequest,
			Code:    "json_binding_error",
			Message: "Invalid JSON format",
			Detail:  err.Error(),
			Context: map[string]interface{}{"error": err.Error()},
		}
	case apperrors.IsNotFound(err):
		return createAPIError(http.StatusNotFound, apperrors.NotFoundError, "Resource not found", err)
	case apperrors.IsValidationError(err):
		return createAPIError(http.StatusBadRequest, apperrors.ValidationError, "Validation error", err)
	case apperrors.IsConflict(err):
		return createAPIError(http.StatusConflict, apperrors.ConflictError, "Resource conflict", err)
	case apperrors.IsAuthorizationError(err):
		return createAPIError(http.StatusForbidden, apperrors.UnauthorizedError, "Unauthorized access", err)
	case apperrors.IsExternalServiceError(err):
		return createAPIError(http.StatusBadGateway, apperrors.ExternalServiceError, "External service failure", err)
	case apperrors.IsSystemError(err):
		return APIError{
			Status:  http.StatusInternalServerError,
			Code:    apperrors.InternalServerError,
			Message: "Internal system error",
			Detail:  "An internal error occurred. Please try again later.",
		}
	default:
		return APIError{
			Status:  http.StatusInternalServerError,
			Code:    "unexpected_error",
			Message: "Unexpected error",
			Detail:  "An unexpected error occurred. Please try again later.",
		}
	}
}

func isJSONError(err error) bool {
	for _, jsonErr := range JSONErrors {
		if errors.As(err, &jsonErr) {
			return true
		}
	}
	return strings.Contains(strings.ToLower(err.Error()), "json:")
}

func createAPIError(status int, code string, message string, err error) APIError {
	apiError := APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Detail:  err.Error(),
	}

	if context, ok := apperrors.GetErrorContext(err); ok {
		apiError.Context = context
	}

	return apiError
}

func formatValidationErrorDetail(ve apperrors.Validation) string {
	switch ve.Rule {
	case apperrors.Required:
		return fmt.Sprintf("The field '%s' is required", ve.Field)
	case apperrors.MinLength:
		minLength, _ := ve.Constraint.(int)
		return fmt.Sprintf("The field '%s' must have at least %d characters", ve.Field, minLength)
	case apperrors.MaxLength:
		maxLength, _ := ve.Constraint.(int)
		return fmt.Sprintf("The field '%s' must have at most %d characters", ve.Field, maxLength)
	case apperrors.InvalidValue:
		return fmt.Sprintf("The value for field '%s' is invalid", ve.Field)
	default:
		return fmt.Sprintf("The field '%s' with value '%v' does not meet the validation rule '%s'",
			ve.Field, ve.Value, ve.Rule)
	}
}

func createValidationErrorContext(ve apperrors.Validation) map[string]interface{} {
	context := map[string]interface{}{
		"field": ve.Field,
		"rule":  string(ve.Rule),
	}

	if ve.Value != nil {
		context["value"] = ve.Value
	}

	if ve.Constraint != nil {
		context["constraint"] = ve.Constraint
	}

	return context
}
