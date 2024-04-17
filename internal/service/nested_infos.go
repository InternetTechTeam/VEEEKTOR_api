package service

import (
	auth "VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetNestedInfosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		NestedInfosGetHandler(w, r)
	case http.MethodPost:
		NestedInfosCreateHandler(w, r)
	case http.MethodPut:
		NestedInfosUpdateHandler(w, r)
	case http.MethodDelete:
		NestedInfosDeleteHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

// Nested infos GET logic.
// If url contains info id, response body will contain all fields.
// If url contains course_id, response will contain array of info pages.
// Info pages can only be accessable for users that belongs to info page course.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// Response: Error message or info pages by course id (info id):
// id : id of info page;
// course_id : id of info page course;
// markdown (optional) : markdown of info page if id is in url values;
// Response codes:
// 200, 400, 401, 403, 404.
func NestedInfosGetHandler(w http.ResponseWriter, r *http.Request) {
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
				w, r, http.StatusNotFound, err)
			return
		}

		belongs, err := models.CheckUserBelongsToCourse(
			claims["user_id"].(int), info.CourseId)
		if err != nil || !belongs {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
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

		belongs, err := models.CheckUserBelongsToCourse(
			claims["user_id"].(int), courseId)
		if err != nil || !belongs {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
			return
		}

		info, err := models.GetNestedInfosByCourseId(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(info)

	} else {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueMissing)
		return
	}

	w.Write(jsonBytes)
}

// Course info POST logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins.
// Response: Error message or StatusOk:
// Expected body:
// course_id : id of course;
// markdown : markdown text of info page.
// Response codes:
// 200, 400, 401, 403.
func NestedInfosCreateHandler(w http.ResponseWriter, r *http.Request) {
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

	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var info models.NestedInfo

	if err = json.Unmarshal(bytes, &info); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), info.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = info.Insert(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested info page PUT logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins.
// Response: Error message or StatusOk:
// Expected body:

// markdown : markdown text of course;
// dep_id : id of course department.
// Response codes:
// 200, 400, 401, 403.
func NestedInfosUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var info models.NestedInfo

	if err = json.Unmarshal(bytes, &info); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), info.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = info.Update(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested info page DELETE logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins that belongs to info course.
// Response: Error message or StatusOk:
// URL values should contain ?id=<id_of_info_page>
// Response codes:
// 200, 400, 401, 403, 404, 500.
func NestedInfosDeleteHandler(w http.ResponseWriter, r *http.Request) {
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
			w, r, http.StatusNotFound, err)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), info.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = models.DeleteNestedInfoById(infoId); err != nil {
		e.ResponseWithError(
			w, r, http.StatusInternalServerError, e.ErrInternalServerError)
		return
	}
}
