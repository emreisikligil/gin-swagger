package api

import "github.com/go-openapi/errors"

type ValidationErrorRespone struct {
	Error string `json:"error"`
	Details []*ValidationError `json:"details,omitempty"`
}

type ValidationError struct {
	Field string `json:"field,omitempty"`
	In string `json:"in,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Msg string `json:"msg,omitempty"`
}

func NewValidationErrorResponse(compositeErr *errors.CompositeError) *ValidationErrorRespone {
	resp := &ValidationErrorRespone{
		Error: "validation_error",
	}

	extractValidationErrors(compositeErr, &resp.Details)

	return resp
}

func extractValidationErrors(compositeErr *errors.CompositeError, validationErrors *[]*ValidationError) {
	for _, err := range compositeErr.Errors {
		if valErr, ok := err.(*errors.Validation); ok {
			*validationErrors = append(*validationErrors, &ValidationError{
				Field: valErr.Name,
				In: valErr.In,
				Value: valErr.Value,
				Msg: valErr.Error(),
			})
		} else if compositeErr, ok := err.(*errors.CompositeError); ok {
			extractValidationErrors(compositeErr, validationErrors)
		}
	}
}