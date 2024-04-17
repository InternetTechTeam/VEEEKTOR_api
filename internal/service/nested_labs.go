package service

import (
	auth "VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetNestedLabsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		NestedLabsGetHandler(w, r)
	case http.MethodPost:
		NestedLabsCreateHandler(w, r)
	case http.MethodPut:
		NestedLabsUpdateHandler(w, r)
	case http.MethodDelete:
		NestedLabsDeleteHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

// Nested labs GET logic.
// If url contains lab id, response body will contain all fields.
// If url contains course_id, response will contain array of labs with few fields.
// Lab pages can only be accessable for users that belongs to lab page course.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// Response: Error message or lab pages by course id (lab id):
// id : lab id;
// course_id : course id;
// opens : date, when lab opens (only on get by id) in UTC;
// closes : date, when lab closes (only on get by id) in UTC;
// topic : lab topic;
// requirements : link to lab requirements (only on get by id);
// example : link to lab example (only on get by id);
// location_id : id of location (only on get by id);
// attempts : number of attempts (only on get by id).
// Response codes:
// 200, 400, 401, 403, 404.
func NestedLabsGetHandler(w http.ResponseWriter, r *http.Request) {
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
		labId, err := strconv.Atoi(rawQuery.Get("id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		lab, err := models.GetNestedLabById(labId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		belongs, err := models.CheckUserBelongsToCourse(
			claims["user_id"].(int), lab.CourseId)
		if err != nil || !belongs {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
			return
		}

		jsonBytes, _ = json.Marshal(lab)

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

		labs, err := models.GetNestedLabsByCourseId(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(labs)

	} else {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueMissing)
		return
	}

	w.Write(jsonBytes)
}

// Nested lab POST logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins.
// Response: Error message or StatusOk:
// Expected body:
// course_id : course id;
// opens : date, when lab opens;
// closes : date, when lab closes;
// topic : lab topic;
// requirements : link to lab requirements (optional);
// example : link to lab example (optional);
// location_id : id of location;
// attempts : number of attempts.
// Response codes:
// 200, 400, 401, 403.
func NestedLabsCreateHandler(w http.ResponseWriter, r *http.Request) {
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

	var lab models.NestedLab

	if err = json.Unmarshal(bytes, &lab); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), lab.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = lab.Insert(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested lab page PUT logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins.
// Response: Error message or StatusOk:
// Expected body:
// id : lab id;
// course_id : course id;
// opens : date, when lab opens;
// closes : date, when lab closes;
// topic : lab topic;
// requirements : link to lab requirements;
// example : link to lab example;
// location_id : id of location;
// attempts : number of attempts.
// Response codes:
// 200, 400, 401, 403.
func NestedLabsUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	var lab models.NestedLab

	if err = json.Unmarshal(bytes, &lab); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), lab.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = lab.Update(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested lab page DELETE logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins that belongs to lab course.
// Response: Error message or StatusOk:
// URL values should contain ?id=<lab_id>
// Response codes:
// 200, 400, 401, 403, 404, 500.
func NestedLabsDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	labId, err := strconv.Atoi(rawQuery.Get("id"))
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
		return
	}

	lab, err := models.GetNestedLabById(labId)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusNotFound, err)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), lab.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = models.DeleteNestedLabById(labId); err != nil {
		e.ResponseWithError(
			w, r, http.StatusInternalServerError, e.ErrInternalServerError)
		return
	}
}
