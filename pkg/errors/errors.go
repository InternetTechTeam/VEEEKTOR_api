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
	ErrOnlyGetAllowed = errors.New(
		"only GET allowed")
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
		"user with this email already exists")
	ErrAccessDenied = errors.New(
		"permission denied")
	// Sessions
	ErrSessionNotExist = errors.New(
		"session for this token doesn't exist")
	ErrTokenExpired = errors.New(
		"token expired")
	ErrTokenNotProvided = errors.New(
		"token not provided")
	ErrTokenNotValid = errors.New(
		"provided token not valid")
	// Roles
	ErrRoleNotFound = errors.New(
		"role with this id not found")
	ErrRoleCantBeSet = errors.New(
		"only admin can set this role")
	// Departments
	ErrDepNotFound = errors.New(
		"department not found")
	ErrDepsNotFound = errors.New(
		"departments not found")
	ErrCantSetThisDep = errors.New(
		"this department can be viewed only by admins")
	// Educational envs
	ErrEdEnvNotFound = errors.New(
		"educational environment not found")
	ErrEdEnvsNotFound = errors.New(
		"educational environments not found")
	// Courses
	ErrCourseIdNull = errors.New(
		"course id must be set")
	ErrCourseNotFound = errors.New(
		"course not found")
	ErrCoursesNotFound = errors.New(
		"courses not found")
	ErrTermNotValid = errors.New(
		"invalid term number")
	ErrTeacherNotFound = errors.New(
		"teacher not found")
	ErrCourseNameNotValid = errors.New(
		"course name not valid")
	ErrUserNotBelongToCourse = errors.New(
		"user not belong to that course")
	// Nested infos
	ErrNestedInfoNotFound = errors.New(
		"nested info page not found")
	ErrNestedInfosNotFound = errors.New(
		"nested info pages not found")
	// Nested labs
	ErrNestedLabNotFound = errors.New(
		"nested lab page not found")
	ErrNestedLabsNotFound = errors.New(
		"nested lab pages not found")
	// Nested tests
	ErrNestedTestNotFound = errors.New(
		"nested test page not found")
	ErrNestedTestsNotFound = errors.New(
		"nested test pages not found")
	ErrTimeLimitTooShort = errors.New(
		"time limit should be more than 5 minutes")
	TimeLimitNotValid = errors.New(
		"time limit not valid")
	// Locations
	ErrLocationNotFound = errors.New(
		"location not found")
)

func ResponseWithError(w http.ResponseWriter, r *http.Request,
	errCode int, err error) {
	w.WriteHeader(errCode)
	w.Write([]byte(fmt.Sprintf(`{"Error" : "%s"}`, err.Error())))
}
