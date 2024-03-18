package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// HTTP
	ErrOnlyPostAllowed = errors.New(
		"only POST allowed")
	ErrUnableToUnmarshalBody = errors.New(
		"unable to unmarshal body")
	ErrMethodNotAllowed = errors.New(
		"method not allowed")
	ErrInternalServerError = errors.New(
		"internal server error")
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
	// Users
	ErrUserNotFound = errors.New(
		"user not found")
	ErrUserExists = errors.New(
		"user with this login already exists")
	// Sessions
	ErrSessionNotExist = errors.New(
		"session for this token doesn't exist")
	ErrTokenExpired = errors.New(
		"token expired")
	ErrTokenNotProvided = errors.New(
		"token not provided")
	// Roles
	ErrRoleNotFound = errors.New(
		"role with this id not found")
	ErrRoleCantBeSet = errors.New(
		"only admin can set this role")
	// Departments
	ErrDepNotFound = errors.New(
		"department with this id not found")
	ErrCantSetThisDep = errors.New(
		"this department can be viewed only by admins")
)

func ResponseWithError(w http.ResponseWriter, r *http.Request,
	errCode int, err error) {
	w.WriteHeader(errCode)
	w.Write([]byte(fmt.Sprintf(`{"Error" : "%s"}`, err.Error())))
}
