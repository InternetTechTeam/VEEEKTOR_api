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
// "refresh_token" : refresh token.
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
// 200, 400, 401, 405, (500).
// If token is expired or session does not exists - code: 401
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
	if exp, err = sess.IsExpired(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusInternalServerError, e.ErrInternalServerError)
		return
	}

	if exp {
		e.ResponseWithError(
			w, r, http.StatusUnauthorized, e.ErrTokenExpired)
		return
	}

	var user models.User
	if user, err = models.GetUserById(sess.UserId); err != nil {
		e.ResponseWithError(
			w, r, http.StatusInternalServerError, err)
		return
	}

	var tokens auth.TokenResponse
	if tokens, err = auth.StoreSession(sess.UserId, user.Id); err != nil {
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

func CheckUserAuthorized(w http.ResponseWriter, r *http.Request) bool {
	accessToken, err := auth.GetAccessTokenFromHeader(r)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusUnauthorized, err)
		return false
	}

	exp, err := auth.IsAccessTokenExpired(accessToken)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusUnauthorized, err)
		return false
	}

	if exp {
		e.ResponseWithError(
			w, r, http.StatusUnauthorized, e.ErrTokenExpired)
		return false
	}

	return true
}
