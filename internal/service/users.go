package service

import (
	"encoding/json"
	"net/http"

	"VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		UsersGetHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrOnlyGetAllowed)
	}
}

// Get user by access token method.
// Expected header:
// Authorization : Bearer <Access token>.
// Response:
// Error message or user data:
// id : user id;
// email : user email;
// name : user name;
// patronymic : user patronymic;
// surname : user surname;
// role_id : id of user role;
// dep_id : id of user department.
// Response codes:
// 200, 400, 401, 404, 500.
func UsersGetHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetAccessTokenFromHeader(r)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	claims, err := auth.GetTokenClaims(accessToken)
	if err == e.ErrTokenExpired {
		e.ResponseWithError(w, r, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := models.GetUserById(claims["user_id"].(int))
	if err != nil {
		e.ResponseWithError(w, r, http.StatusNotFound, e.ErrUserNotFound)
		return
	}
	user.Password = ""

	jsonBytes, _ := json.Marshal(user)
	w.Write(jsonBytes)
}

// Users authorization, authentication logic.
// Expected body:
// email : user email (4-64 symbols);
// password : user password (8-50 symbols).
// Response:
// Error message or token pair:
// access_token  : token for access to private pages, lifetime - 15m;
// refresh_token : token for refreshing access token, lifetime - 30 days.
// Access token claims:
// exp : token expiration date and time in UNIX format;
// user_id : user id;
// role_id : user role id.
// Cookie:
// refresh_token : <rt>.
// Response codes:
// 200, 400, 404, 405.
func UsersSignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ResponseWithError(
			w, r, http.StatusMethodNotAllowed, e.ErrOnlyPostAllowed)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var err error
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

	tokens, _ := auth.StoreSession(user.Id, user.RoleId, user.GroupId)

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

// Users sign up logic.
// Expected body:
// email : user email (4-64 symbols);
// password : user password (8-50 symbols);
// name : user name (2-30 symbols);
// patronymic : user patronymic (2-30 symbols);
// surname : user surname (2-30 symbols);
// dep_id : department id.
// Response:
// Error message or StatusOk.
// Response codes:
// 200, 400, 405, 409.
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
	dto.RoleId = 1

	if err = json.Unmarshal(bytes, &dto); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	// Teachers and admins can not be created via this function
	if dto.RoleId != 1 {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrRoleCantBeSet)
		return
	}

	// Admin department not availiable for basic users
	if dto.DepId == 1 {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrCantSetThisDep)
	}

	if err := dto.Validate(); err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := dto.Insert(); err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
