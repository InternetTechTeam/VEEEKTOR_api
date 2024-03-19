package service

import (
	"VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
)

func GetCouresesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		CoursesGetByUserIdHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrOnlyGetAllowed)
	}
}

func CoursesGetByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	if !CheckUserAuthorized(w, r) {
		return
	}

	accessToken, err := auth.GetAccessTokenFromHeader(r)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusUnauthorized,
			err)
		return
	}

	claims, err := auth.GetTokenClaims(accessToken)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusUnauthorized,
			err)
		return
	}

	userId := int(claims["user_id"].(float64))

	courses, err := models.GetCoursesByUserId(userId)
	if err != nil {
		e.ResponseWithError(w, r, http.StatusNotFound, err)
	}

	jsonBytes, _ := json.Marshal(courses)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
