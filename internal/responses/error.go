// Brief   : Structure used for error logging back to the client.
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 17.06.2023

//nolint:gomnd
package res

import "net/http"

// All errors thrown
// General errors
var (
	ErrNotFound   = NewError(http.StatusNotFound, 0000, "Resource not found", "Resource not found")
	ErrDB         = NewError(http.StatusInternalServerError, 1010, "Server error", "Database error")
	ErrForbidden  = NewError(http.StatusForbidden, 0000, "Not allowed to read requested resource", "")
	ErrBadRequest = NewError(http.StatusBadRequest, 0000, "The request body is invalid", "")
)

// JWT Claim
var (
	ErrMissingAuth    = NewError(http.StatusUnauthorized, 1000, "Unauthenticated", "Missing auth token")
	ErrBadAuthHeader  = NewError(http.StatusUnauthorized, 1001, "Unauthenticated", "Bad auth header format")
	ErrBadToken       = NewError(http.StatusUnauthorized, 1002, "Unauthenticated", "Unable to parse token")
	ErrMissingSubject = NewError(http.StatusUnauthorized, 1003, "Unauthenticated", "Token has no subject")
)

// Team errors
var (
	ErrParticipantsExist = NewError(http.StatusBadRequest, 1000, "Some atheltes are already part of the team", "")
)

type Error struct {
	Error         int    `json:"error" validate:"required"`
	ErrorSubCode  int    `json:"-"`
	Message       string `json:"message" validate:"required"`
	HiddenMessage string `json:"-"`
	InternalMsg   string `json:"-"`
}

// NewError creates a new error instance with a code and a human-readable message.
func NewError(errorCode int, errorSubCode int, message string, hiddenMessage string) *Error {
	return &Error{
		errorCode,
		errorSubCode,
		message,
		hiddenMessage,
		"",
	}
}

// Err attaches an internal error message to the error
func (e *Error) Err(internalError error) *Error {
	e.InternalMsg = internalError.Error()
	return e
}
