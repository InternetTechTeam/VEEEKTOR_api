package service

import (
	auth "VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GroupsGetHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

// Groups GET logic.
// Groups can be get via group id or department id.
// Url values should contain ?id=<group_id> or ?dep_id=<department_id>
// Response: Error message or group:
// id : group id
// name : group name like O722B
// dep_id : group department id
// Response codes:
// 200, 400, 404.
func GroupsGetHandler(w http.ResponseWriter, r *http.Request) {
	var jsonBytes []byte

	rawQuery := r.URL.Query()
	if rawQuery.Has("id") {
		groupId, err := strconv.Atoi(rawQuery.Get("id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		group, err := models.GetGroupById(groupId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(group)

	} else if rawQuery.Has("dep_id") {
		depId, err := strconv.Atoi(rawQuery.Get("dep_id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		groups, err := models.GetAllGroupsByDepId(depId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, e.ErrGroupsNotFound)
			return
		}

		jsonBytes, _ = json.Marshal(groups)
	} else {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueMissing)
		return
	}

	w.Write(jsonBytes)
}

// Groups linkage to course logic.
// Expected header:
// Authorization : Bearer <access token>.
// This method allowed only to teachers, who belongs to course or admins.
// Response: Error message or StatusOk:
// Expected body:
// group_id : id of group;
// course_id : id of course.
// Response codes:
// 200, 400, 401, 403.
func LinkGroupWithCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ResponseWithError(
			w, r, http.StatusMethodNotAllowed, e.ErrOnlyPostAllowed)
		return
	}

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

	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var gc models.GroupCourse
	if err := json.Unmarshal(bytes, &gc); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	var course models.Course
	course.Id = gc.CourseId
	access, err := course.CheckAccess(claims)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrCourseNotFound)
		return
	}
	if access != 2 {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = gc.Insert(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Groups unlink from course logic.
// Expected header:
// Authorization : Bearer <access token>.
// This method allowed only to teachers, who belongs to course or admins.
// Response: Error message or StatusOk:
// Expected body:
// group_id : id of group;
// course_id : id of course.
// Response codes:
// 200, 400, 401, 403.
func UnlinkGroupFromCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ResponseWithError(
			w, r, http.StatusMethodNotAllowed, e.ErrOnlyPostAllowed)
		return
	}

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

	if claims["role_id"] != 2 && claims["role_id"] != 3 {
		e.ResponseWithError(w, r, http.StatusForbidden, e.ErrAccessDenied)
		return
	}

	bytes := make([]byte, r.ContentLength)
	r.Body.Read(bytes)

	var gc models.GroupCourse

	if err := json.Unmarshal(bytes, &gc); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	var course models.Course
	course.Id = gc.CourseId
	access, err := course.CheckAccess(claims)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest,
			e.ErrCourseNotFound)
		return
	}
	if access != 2 {
		e.ResponseWithError(w, r, http.StatusForbidden,
			e.ErrUserNotBelongToCourse)
		return
	}

	if err = gc.Delete(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
