package service

import (
	auth "VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetNestedTestsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		NestedTestsGetHandler(w, r)
	case http.MethodPost:
		NestedTestsCreateHandler(w, r)
	case http.MethodPut:
		NestedTestsUpdateHandler(w, r)
	case http.MethodDelete:
		NestedTestsDeleteHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

// Nested tests GET logic.
// If url contains test id, response body will contain all fields.
// If url contains course_id, response will contain array of tests with few fields.
// Test pages can only be accessable for users that belongs to test page course.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// Response: Error message or info pages by course id (info id):
// id : test id;
// course_id : course id;
// opens : date, when test opens (only on get by id) in UTC;
// closes : date, when test closes (only on get by id) in UTC;
// tasks_count : count of test tasks (only on get by id);
// topic : test topic;
// location_id : id of location (only on get by id);
// attempts : number of attempts (only on get by id);
// password : test password (optional) (only on get by id);
// time_limit : time limit duration (only on get by id).
// Response codes:
// 200, 400, 401, 403, 404.
func NestedTestsGetHandler(w http.ResponseWriter, r *http.Request) {
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
		testId, err := strconv.Atoi(rawQuery.Get("id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		test, err := models.GetNestedTestById(testId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		belongs, err := models.CheckUserBelongsToCourse(
			claims["user_id"].(int), test.CourseId)
		if err != nil || !belongs {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
			return
		}

		jsonBytes, _ = json.Marshal(test)

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

		tests, err := models.GetNestedTestsByCourseId(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(tests)

	} else {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueMissing)
		return
	}

	w.Write(jsonBytes)
}

// Nested test POST logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins.
// Response: Error message or StatusOk:
// Expected body:
// course_id : course id;
// opens : date, when test opens in UTC;
// closes : date, when test closes in UTC;
// tasks_count : count of test tasks;
// topic : test topic;
// location_id : id of location;
// attempts : number of attempts;
// password : test password (optional);
// time_limit : time limit duration (00:15:00).
// Response codes:
// 200, 400, 401, 403.
func NestedTestsCreateHandler(w http.ResponseWriter, r *http.Request) {
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

	var test models.NestedTest

	if err = json.Unmarshal(bytes, &test); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), test.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = test.Insert(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested test page PUT logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins.
// Response: Error message or StatusOk:
// Expected body:
// id : test id;
// course_id : course id;
// opens : date, when test opens in UTC;
// closes : date, when test closes in UTC;
// tasks_count : count of test tasks;
// topic : test topic;
// location_id : id of location;
// attempts : number of attempts;
// password : test password (optional);
// time_limit : time limit duration (00:15:00).
// Response codes:
// 200, 400, 401, 403.
func NestedTestsUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	var test models.NestedTest

	if err = json.Unmarshal(bytes, &test); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUnableToUnmarshalBody)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), test.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = test.Update(); err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Nested test page DELETE logic.
// Expected header:
// Authorization : Bearer <Valid Access Token>
// This method allowed only to teachers and admins that belongs to test course.
// Response: Error message or StatusOk:
// URL values should contain ?id=<test_id>
// Response codes:
// 200, 400, 401, 403, 404, 500.
func NestedTestsDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	testId, err := strconv.Atoi(rawQuery.Get("id"))
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
		return
	}

	test, err := models.GetNestedTestById(testId)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusNotFound, err)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), test.CourseId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	if err = models.DeleteNestedTestById(testId); err != nil {
		e.ResponseWithError(
			w, r, http.StatusInternalServerError, e.ErrInternalServerError)
		return
	}
}
