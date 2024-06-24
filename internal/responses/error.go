// Brief   : Structure used for error logging back to the client.
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 17.06.2023

//nolint:mnd
package res

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

// All errors thrown
// General errors
var (
	ErrNotFound          = NewError(http.StatusNotFound, 0000, "Resource not found", "Resource not found")
	ErrDB                = NewError(http.StatusInternalServerError, 0010, "Server error", "Database error")
	ErrZitadel           = NewError(http.StatusInternalServerError, 0020, "Server error", "Error communicating with Zitadel")
	ErrForbidden         = NewError(http.StatusForbidden, 0000, "Not allowed to read requested resource", "")
	ErrForbiddenNotFound = NewError(http.StatusForbidden, 0010, "Not allowed to read requested resource", "The resource does not exist")
	ErrBadRequest        = NewError(http.StatusBadRequest, 0000, "The request body is invalid", "")
	ErrMissingTotp       = NewError(http.StatusUnauthorized, 0010, "2FA must be re-verified", "")
)

// JWT Claim
var (
	ErrMissingAuth    = NewError(http.StatusUnauthorized, 0000, "Unauthenticated", "Missing auth token")
	ErrBadAuthHeader  = NewError(http.StatusUnauthorized, 0001, "Unauthenticated", "Bad auth header format")
	ErrBadToken       = NewError(http.StatusUnauthorized, 0002, "Unauthenticated", "Unable to parse token")
	ErrMissingSubject = NewError(http.StatusUnauthorized, 0003, "Unauthenticated", "Token has no subject")
)

// Team errors
var (
	ErrParticipantsExist = NewError(http.StatusBadRequest, 1000, "Some atheltes are already part of the team", "")
)

// Event errors
var (
	ErrEventHasResults = NewError(http.StatusBadRequest, 2000, "Results are registered for this event.", "")
)

type Error struct {
	Error         int    `json:"error" validate:"required"`
	ErrorSubCode  int    `json:"subCode"`
	Message       string `json:"message" validate:"required"`
	HiddenMessage string `json:"hiddenMessage"`
	InternalMsg   string `json:"internalMessage"`
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

func NewErrorFromDb(err error) *Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound.Err(err)
	} else {
		return ErrDB.Err(err)
	}
}

// Err attaches an internal error message to the error
func (e *Error) Err(internalError error) *Error {
	e.InternalMsg = internalError.Error()
	return e
}

func (e *Error) String() string {
	return fmt.Sprintf("%v-%v: %v, Hidden Message: %v, Internal Message: %v", e.Error, e.ErrorSubCode, e.Message, e.HiddenMessage, e.InternalMsg)
}
