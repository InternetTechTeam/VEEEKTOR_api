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
// Lab pages can only be accessible for users that belongs to course.
// Url values should contain ?id=<test_id> or ?id=<course_id>.
// Expected header:
// Authorization : Bearer <access token>
// Response: Error message or test(s) by test id (course id):
// id : test id;
// course_id : course id;
// opens : date, when test opens in UTC;
// closes : date, when test closes in UTC;
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
		e.ResponseWithError(w, r, http.StatusUnauthorized, e.ErrTokenExpired)
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
				w, r, http.StatusNotFound, e.ErrNestedTestNotFound)
			return
		}

		if belongs, _ := models.CheckUserBelongsToCourse(
			claims["user_id"].(int), test.CourseId); !belongs {
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

		if belongs, _ := models.CheckUserBelongsToCourse(
			claims["user_id"].(int), courseId); !belongs {
			e.ResponseWithError(
				w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
			return
		}

		tests, err := models.GetNestedTestsByCourseId(courseId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, e.ErrNestedTestsNotFound)
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
// Authorization : Bearer <access token>
// This method allowed only to teachers, who belongs to course or admins.
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
		e.ResponseWithError(w, r, http.StatusUnauthorized, e.ErrTokenExpired)
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

	if belongs, _ := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), test.CourseId); !belongs {
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
// Authorization : Bearer <access token>
// This method allowed only to teachers, who belongs to course or admins.
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
		e.ResponseWithError(w, r, http.StatusUnauthorized, e.ErrTokenExpired)
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

	if belongs, _ := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), test.CourseId); !belongs {
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
// URL values should contain ?id=<test_id>
// Expected header:
// Authorization : Bearer <access token>
// This method allowed only to admins and teachers, who belongs to course.
// Response: Error message or StatusOk:
// Response codes:
// 200, 400, 401, 403, 404.
func NestedTestsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetAccessTokenFromHeader(r)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusBadRequest, err)
		return
	}

	claims, err := auth.GetTokenClaims(accessToken)
	if err == e.ErrTokenExpired {
		e.ResponseWithError(w, r, http.StatusUnauthorized, e.ErrTokenExpired)
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
			w, r, http.StatusNotFound, e.ErrNestedTestNotFound)
		return
	}

	if belongs, _ := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), test.CourseId); !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	_ = models.DeleteNestedTestById(testId)

	w.WriteHeader(http.StatusOK)
}
