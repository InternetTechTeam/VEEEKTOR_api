package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Authorization, authentication logic.
// Expected body:
// "email" 	  : user email (4-64 symbols),
// "password" : user password (8-50 symbols).
// Response:
// Error message or token pair:
// "access_token"  : token for access to private pages, lifetime - one hour,
// "refresh_token" : token for refreshing access token, lifetime - 30 days.
// Access token claims:
// "exp" : token expiration date and time,
// "user_id" : ID of the user who owns the token,
// "role_id" : User role id. For actual roles see roles API.
// Cookie:
// "refresh_token".
// Response codes:
// 200, 400, 404, 405, (500).
func UsersSignInHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodPost {
		e.ResponseWithError(
			w, r, http.StatusMethodNotAllowed, e.ErrOnlyPostAllowed)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var inp models.SignInInput
	if err = json.Unmarshal(bytes, &inp); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	if err := inp.Validate(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	var user models.User
	if user, err = models.GetUserByEmailAndPassword(inp); err != nil {
		e.ResponseWithError(
			w, r, http.StatusNotFound, err)
		return
	}

	var tokens auth.TokenResponse
	if tokens, err = auth.StoreSession(user.Id, user.RoleId); err != nil {
		e.ResponseWithError(
			w, r, http.StatusInternalServerError, err)
		return
	}

	// Write jwt and refresh token pair
	jsonBytes, _ := json.Marshal(tokens)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Write(jsonBytes)
}

// User creation logic.
// Expected body:
// "email" 	  : user email (4-64 symbols),
// "password" : user password (8-50 symbols),
// "name" : user name (2-30 symbols),
// "patronymic" : user patronymic (2-30 symbols),
// "surname" : user surname (2-30 symbols),
// "role_id" : ID of user role.
// Response:
// Error message or null.
// Response codes:
// 200, 400, 405, 409, (500).
func UsersSignUpHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodPost {
		e.ResponseWithError(
			w, r, http.StatusMethodNotAllowed, e.ErrOnlyPostAllowed)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var dto models.User
	if err = json.Unmarshal(bytes, &dto); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	if err := dto.Validate(); err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := dto.Insert(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				e.ResponseWithError(
					w, r, http.StatusConflict, e.ErrUserExists)
				return
			}
		} else {
			e.ResponseWithError(
				w, r, http.StatusInternalServerError,
				e.ErrInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
