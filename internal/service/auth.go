package service

import (
	"encoding/json"
	"net/http"

	auth "VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
)

// Token refresh for mobile and web clients.
// Expected cookie / body (for mobile clients):
// "refresh_token" : <refresh token>.
// Response:
// Error message or token pair:
// "access_token"  : token for access to private pages, lifetime - 15m;
// "refresh_token" : token for refreshing access token, lifetime - 30 days.
// Access token claims:
// "exp" : token expiration date and time;
// "user_id" : user id;
// "role_id" : user role id.
// Cookie:
// "refresh_token" : <rt>.
// Response codes:
// 200, 400, 401, 405.
func UpdateToken(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodPost {
		e.ResponseWithError(
			w, r, http.StatusMethodNotAllowed, e.ErrOnlyPostAllowed)
		return
	}

	var refreshToken string
	if refreshToken, err = auth.GetRefreshTokenFromCookieOrBody(r); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	var sess auth.Session
	if sess, err = auth.GetSessionByRefreshToken(refreshToken); err != nil {
		e.ResponseWithError(
			w, r, http.StatusUnauthorized, err)
		return
	}

	var exp bool
	if exp, _ = sess.IsExpired(); exp {
		e.ResponseWithError(
			w, r, http.StatusUnauthorized, e.ErrTokenExpired)
		return
	}

	var user models.User
	if user, err = models.GetUserById(sess.UserId); err != nil {
		e.ResponseWithError(
			w, r, http.StatusNotFound, e.ErrUserNotFound)
		return
	}

	tokens, _ := auth.StoreSession(sess.UserId, user.Id)

	// Write jwt and refresh token pair
	jsonBytes, _ := json.Marshal(tokens)
	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: tokens.RefreshToken,
		Path:  "/",
		// Secure:   true, *causes cookie not set with http unsecure protocol*
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Write(jsonBytes)
}

// Log out logic.
// Expected cookie / body (for mobile clients):
// "refresh_token" : <refresh token>.
// Response:
// Error message or StatusOk.
// Cookie:
// "refresh_token" : <null>.
// Response codes:
// 200, 400, 405.
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ResponseWithError(
			w, r, http.StatusMethodNotAllowed, e.ErrOnlyPostAllowed)
		return
	}

	var err error
	var refreshToken string
	if refreshToken, err = auth.GetRefreshTokenFromCookieOrBody(r); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrTokenNotProvided)
		return
	}

	_ = auth.DeleteSessionByRT(refreshToken)

	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: "",
		Path:  "/",
		// Secure:   true, *causes cookie not set with http unsecure protocol*
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
