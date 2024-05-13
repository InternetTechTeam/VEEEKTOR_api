package service

import (
	"VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
)

func GetCouresesHandler(w http.ResponseWriter, r *http.Request) {
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
		CoursesGetHandler(w, r, token, claims)
	case http.MethodPost:
		CoursesCreateHandler(w, r, token, claims)
	case http.MethodPut:
		CoursesUpdateHandler(w, r, token, claims)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

// Courses GET logic.
// Courses can be get by id in url values or by group_id in token claims;
// Expected header:
// Authorization : Bearer <access token>
// Response: Error message or course(s) by course_id (group_id):
// id : id of course;
// name : name of course;
// term : term of course;
// teacher_id : id of teacher (user) (get by course_id only);
// markdown : markdown text of course (get by course_id only);
// dep_id : id of course department (get by course_id only);
// teacher.name : teacher name (get by group_id only);
// teacher.patronymic : teacher patronymic (get by group_id only);
// teacher.surname : teacher surname (get by group_id only);
// teacher.dep : teacher department (get by group_id only);
// dep : course department (get by group_id only);
// modified_at : course last modified time in UNIX format (get by group_id only).
// Response codes:
// 200, 400, 401, 403, 404.
func CoursesGetHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	var err error
	var jsonBytes []byte

	rawQuery := r.URL.Query()
	if rawQuery.Has("id") {
		courseId, err := strconv.Atoi(rawQuery.Get("id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		course, err := models.GetCourseById(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, e.ErrCourseNotFound)
			return
		}

		// Error checked in GetCourseById
		access, _ := course.CheckAccess(claims)
		if access == 0 {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrAccessDenied)
			return
		}

		jsonBytes, _ = json.Marshal(course)

	} else {
		var courses []models.CourseMultipleExportDTO
		if claims["role_id"] == 1 { // Student
			courses, err = models.GetAllCoursesByGroupId(claims["group_id"].(int))
			if err != nil {
				e.ResponseWithError(w, r, http.StatusNotFound, e.ErrCoursesNotFound)
				return
			}
		} else { // Teacher or admin
			courses, err = models.GetAllCoursesByTeacherId(claims["user_id"].(int))
			if err != nil {
				e.ResponseWithError(w, r, http.StatusNotFound, e.ErrCoursesNotFound)
				return
			}
		}

		jsonBytes, _ = json.Marshal(courses)
	}

	w.Write(jsonBytes)
}

// Courses INSERT logic.
// Expected header:
// Authorization : Bearer <access token>
// Course creation allowed only to teachers and admins.
// Response:
// id : course id.
// Expected body:
// name : name of course;
// term : term of course;
// teacher_id : id of teacher (user);
// markdown : markdown text of course;
// dep_id : id of course department.
// Response codes:
// 200, 400, 401, 403.
func CoursesCreateHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var course models.Course

	if err := json.Unmarshal(bytes, &course); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	course_id, err := course.Insert()
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"id" : %d}`, course_id)))
}

// Courses PUT logic.
// Expected header:
// Authorization : Bearer <access token>
// Course update allowed only to teachers and admins.
// Response: Error message or StatusOk:
// Expected body:
// id : id of course;
// name : name of course;
// term : term of course;
// teacher_id : id of teacher (user);
// markdown : markdown text of course;
// dep_id : id of course department.
// Response codes:
// 200, 400, 401, 403.
func CoursesUpdateHandler(w http.ResponseWriter, r *http.Request,
	token string, claims jwt.MapClaims) {
	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var course models.Course

	if err := json.Unmarshal(bytes, &course); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	access, err := course.CheckAccess(claims)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrCourseNotFound)
		return
	}
	if access != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	if err := course.Update(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
