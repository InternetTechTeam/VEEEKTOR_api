package service

import (
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

	} else if rawQuery.Has("env_id") {
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
	}

	w.Write(jsonBytes)
}
