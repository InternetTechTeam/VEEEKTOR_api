package service

import (
	auth "VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
)

func GetNestedInfosHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAccessTokenFromHeader(r)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	claims, err := auth.GetTokenClaims(token)
	if err == e.ErrTokenExpired {
		e.ResponseWithError(w, r, http.StatusUnauthorized, e.ErrTokenExpired)
		return
	} else if err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		NestedInfosGetHandler(w, r, token, claims)
	case http.MethodPost:
		NestedInfosCreateHandler(w, r, token, claims)
	case http.MethodPut:
		NestedInfosUpdateHandler(w, r, token, claims)
	case http.MethodDelete:
		NestedInfosDeleteHandler(w, r, token, claims)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

// Nested infos GET logic.
// Url values should contain ?id=<info_id> or ?id=<course_id>.
// Info pages can only be accessible for users that belongs to info page course.
// Expected header:
// Authorization : Bearer <access token>
// Response: Error message or info page(s) by info id (course id):
// id : id of info page;
// course_id : id of info page course;
// markdown : markdown of info page (only for get by id);
// Response codes:
// 200, 400, 401, 403, 404.
func NestedInfosGetHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	var jsonBytes []byte

	rawQuery := r.URL.Query()
	if rawQuery.Has("id") {
		infoId, err := strconv.Atoi(rawQuery.Get("id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		info, err := models.GetNestedInfoById(infoId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, e.ErrNestedInfoNotFound)
			return
		}

		if info.CheckAccess(claims) == 0 {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrAccessDenied)
			return
		}

		jsonBytes, _ = json.Marshal(info)

	} else if rawQuery.Has("course_id") {
		courseId, err := strconv.Atoi(rawQuery.Get("course_id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		var course models.Course
		course.Id = courseId

		if course.CheckAccess(claims) == 0 {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
			return
		}

		infos, err := models.GetNestedInfosByCourseId(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, e.ErrNestedInfosNotFound)
			return
		}

		jsonBytes, _ = json.Marshal(infos)

	} else {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueMissing)
		return
	}

	w.Write(jsonBytes)
}

// Course infos POST logic.
// Expected header:
// Authorization : Bearer <access token>
// This method allowed only to teachers, who belongs to course or admins.
// Response: Error message or StatusOk:
// Expected body:
// course_id : id of course;
// markdown : markdown text of info page.
// Response codes:
// 200, 400, 401, 403.
func NestedInfosCreateHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var info models.NestedInfo

	if err := json.Unmarshal(bytes, &info); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	if info.CheckAccess(claims) != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err := info.Insert(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested infos PUT logic.
// Expected header:
// Authorization : Bearer <access token>
// This method allowed only to teachers, who belongs to course or admins.
// Response: Error message or StatusOk:
// Expected body:
// id : nested info page id;
// course_id : nested info page course id;
// name : nested info page name;
// markdown : markdown of nested info page.
// Response codes:
// 200, 400, 401, 403.
func NestedInfosUpdateHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var info models.NestedInfo

	if err := json.Unmarshal(bytes, &info); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	if info.CheckAccess(claims) != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	if err := info.Update(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested infos DELETE logic.
// URL values should contain ?id=<info_page_id>.
// Expected header:
// Authorization : Bearer <access token>.
// This method allowed only to admins and teachers, who belongs to course.
// Response: Error message or StatusOk:
// Response codes:
// 200, 400, 401, 403, 404.
func NestedInfosDeleteHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	rawQuery := r.URL.Query()
	if !rawQuery.Has("id") {
		e.ResponseWithError(w, r, http.StatusBadRequest, e.ErrUrlValueMissing)
		return
	}

	infoId, err := strconv.Atoi(rawQuery.Get("id"))
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
		return
	}

	info, err := models.GetNestedInfoById(infoId)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusNotFound, e.ErrNestedInfoNotFound)
		return
	}

	if info.CheckAccess(claims) != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	_ = models.DeleteNestedInfoById(infoId)

	w.WriteHeader(http.StatusOK)
}
