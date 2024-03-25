package service

import (
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetEducatinalEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		EducationalEnvsGetByUserIdHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrOnlyGetAllowed)
	}
}

// Get all educational environments logic. If url
// contains id, response will contain educational env by id.
// Response: Error message or educational environment(s):
// id : id of educational env;
// name : name of educational env.
// Response codes:
// 200, 400, 404, 500.
func EducationalEnvsGetByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	var jsonBytes []byte

	rawQuery := r.URL.Query()
	if rawQuery.Has("id") {
		envId, err := strconv.Atoi(rawQuery.Get("id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		env, err := models.GetEducationalEnvironmentById(envId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(env)
	} else {
		envs, err := models.GetAllEducationalEnvs()
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusInternalServerError, err)
			return
		}

		jsonBytes, _ = json.Marshal(envs)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
