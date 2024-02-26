package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrOnlyPostAllowed = errors.New(
		"only POST allowed")
	ErrUnableToUnmarshalBody = errors.New(
		"unable to unmarshal body")
	ErrInternalServerError = errors.New(
		"internal server error")
	ErrUserNotFound = errors.New(
		"user not found")
	ErrUserExists = errors.New(
		"user with this login already exists")
	ErrSessionNotExist = errors.New(
		"session for this token doesn't exist")
	ErrMethodNotAllowed = errors.New(
		"method not allowed")
	ErrUrlValueNotValid = errors.New(
		"url value is not valid")
	ErrFieldViolatesFK = errors.New(
		"fields violates foreign key")
	ErrUrlValueMissing = errors.New(
		"url value is missing")
	ErrMissingFields = errors.New(
		"required fields are missing")
	ErrCantPrepareDbStmt = errors.New(
		"cant prepare db statement")
	ErrRoleNotFound = errors.New(
		"role with this id not found")
	ErrTokenExpired = errors.New(
		"token expired")
	ErrTokenNotProvided = errors.New(
		"token not provided")
)

func ResponseWithError(w http.ResponseWriter, r *http.Request,
	errCode int, err error) {
	w.WriteHeader(errCode)
	w.Write([]byte(fmt.Sprintf(`{"Error" : "%s"}`, err.Error())))
}
