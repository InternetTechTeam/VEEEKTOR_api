package service

import (
	auth "VEEEKTOR_api/internal/auth"
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetNestedInfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		NestedInfoGetHandler(w, r)
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrMethodNotAllowed)
	}
}

func NestedInfoGetHandler(w http.ResponseWriter, r *http.Request) {
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

	rawQuery := r.URL.Query()
	if !rawQuery.Has("id") {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueMissing)
		return
	}

	infoId, err := strconv.Atoi(rawQuery.Get("id"))
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
		return
	}

	belongs, err := models.CheckUserBelongsToCourse(
		claims["user_id"].(int), infoId)
	if err != nil || !belongs {
		e.ResponseWithError(
			w, r, http.StatusForbidden, e.ErrUserNotBelongToCourse)
		return
	}

	info, err := models.GetNestedInfoById(infoId)
	if err != nil {
		e.ResponseWithError(
			w, r, http.StatusBadRequest, err)
		return
	}

	var jsonBytes []byte
	jsonBytes, _ = json.Marshal(info)

	w.Write(jsonBytes)
}
