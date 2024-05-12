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

func GetNestedLabsHandler(w http.ResponseWriter, r *http.Request) {
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
		NestedLabsGetHandler(w, r, token, claims)
	case http.MethodPost:
		NestedLabsCreateHandler(w, r, token, claims)
	case http.MethodPut:
		NestedLabsUpdateHandler(w, r, token, claims)
	case http.MethodDelete:
		NestedLabsDeleteHandler(w, r, token, claims)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

// Nested labs GET logic.
// Lab pages can only be accessible for users that belongs to course.
// Url values should contain ?id=<lab_id> or ?id=<course_id>.
// Expected header:
// Authorization : Bearer <access token>
// Response: Error message or lab page(s) by lab id (course id):
// id : lab id;
// course_id : course id;
// opens : date, when lab opens in UTC;
// closes : date, when lab closes in UTC;
// topic : lab topic;
// requirements : link to lab requirements (only with get by id);
// example : link to lab example (only with get by id);
// location_id : id of location (only with get by id);
// attempts : number of attempts (only with get by id).
// Response codes:
// 200, 400, 401, 403, 404.
func NestedLabsGetHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
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
				w, r, http.StatusNotFound, e.ErrNestedLabNotFound)
			return
		}

		if lab.CheckAccess(claims) == 0 {
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

		var course models.Course
		course.Id = courseId

		if course.CheckAccess(claims) == 0 {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
			return
		}

		labs, err := models.GetNestedLabsByCourseId(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, e.ErrNestedLabsNotFound)
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

// Nested labs POST logic.
// Expected header:
// Authorization : Bearer <access token>
// This method allowed only to teachers, who belongs to course or admins.
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
func NestedLabsCreateHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var lab models.NestedLab

	if err := json.Unmarshal(bytes, &lab); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	if lab.CheckAccess(claims) != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	if err := lab.Insert(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested labs PUT logic.
// Expected header:
// Authorization : Bearer <access token>
// This method allowed only to teachers, who belongs to course or admins.
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
func NestedLabsUpdateHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var lab models.NestedLab
	if err := json.Unmarshal(bytes, &lab); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	if lab.CheckAccess(claims) != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	if err := lab.Update(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested labs DELETE logic.
// URL values should contain ?id=<lab_id>
// Expected header:
// Authorization : Bearer <access token>
// This method allowed only to admins and teachers, who belongs to course.
// Response: Error message or StatusOk:
// Response codes:
// 200, 400, 401, 403, 404, 500.
func NestedLabsDeleteHandler(w http.ResponseWriter, r *http.Request,
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

	labId, err := strconv.Atoi(rawQuery.Get("id"))
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
		return
	}

	lab, err := models.GetNestedLabById(labId)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusNotFound, e.ErrNestedLabNotFound)
		return
	}

	if lab.CheckAccess(claims) != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	_ = models.DeleteNestedLabById(labId)

	w.WriteHeader(http.StatusOK)
}
