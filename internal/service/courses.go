package service

import (
	"VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetCouresesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		CoursesGetByUserIdHandler(w, r)
	case http.MethodPost:
		CoursesCreateHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrOnlyGetAllowed)
	}
}

// Get all courses without markdown by user id logic. If url
// contains id, response body will contain markdown.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// Response: Error message or courses by user id (course id):
// id : id of course;
// name : name of course;
// term : term of course;
// teacher_id : id of teacher (user);
// markdown : markdown text of course;
// dep_id : id of course department.
// Response codes:
// 200, 400, 401, 404.
func CoursesGetByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	authorized, err := auth.CheckUserAuthorized(r)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusUnauthorized, err)
		return
	}

	if !authorized {
		e.ResponseWithError(w, r, http.StatusUnauthorized, e.ErrTokenExpired)
		return
	}

	var jsonBytes []byte

	rawQuery := r.URL.Query()
	if rawQuery.Has("id") {

		var courseId int
		if courseId, err = strconv.Atoi(rawQuery.Get("id")); err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		course, err := models.GetCourseById(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(course)

	} else {
		userId, err := auth.GetUserIdFromRequest(r)
		if err != nil {
			e.ResponseWithError(w, r, http.StatusBadRequest, err)
			return
		}

		courses, err := models.GetAllCoursesByUserId(userId)
		if err != nil {
			e.ResponseWithError(w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(courses)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Course insert logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// Response: Error message StatusOk:
// Expected body:
// name : name of course;
// term : term of course;
// teacher_id : id of teacher (user);
// markdown : markdown text of course;
// dep_id : id of course department.
// Response codes:
// 200, 400, 401.
func CoursesCreateHandler(w http.ResponseWriter, r *http.Request) {
	authorized, err := auth.CheckUserAuthorized(r)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusUnauthorized, err)
		return
	}

	if !authorized {
		e.ResponseWithError(w, r, http.StatusUnauthorized, e.ErrTokenExpired)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var course models.Course

	if err = json.Unmarshal(bytes, &course); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	if err = course.Insert(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
