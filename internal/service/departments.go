package service

import (
	"VEEEKTOR_api/internal/models"
	e "VEEEKTOR_api/pkg/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		DepartmentsGetHandler(w, r)
	default:
		e.ResponseWithError(w, r, http.StatusMethodNotAllowed,
			e.ErrOnlyGetAllowed)
	}
}

// Get all departments logic. If url
// contains id, response will contain department by id.
// If url contains env_id, response will contain departments
// by environment id.
// Response: Error message or department(s):
// id : id of department;
// name : name of department;
// env_id : id of department educational environment.
// Response codes:
// 200, 400, 404, 500.
func DepartmentsGetHandler(w http.ResponseWriter, r *http.Request) {
	var jsonBytes []byte

	rawQuery := r.URL.Query()
	if rawQuery.Has("id") {
		depId, err := strconv.Atoi(rawQuery.Get("id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		dep, err := models.GetDepartmentById(depId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusNotFound, err)
			return
		}

		jsonBytes, _ = json.Marshal(dep)
	} else if rawQuery.Has("env_id") {
		envId, err := strconv.Atoi(rawQuery.Get("env_id"))
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusBadRequest, e.ErrUrlValueNotValid)
			return
		}

		deps, err := models.GetAllDepartmentsByEnvironmentId(envId)
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusInternalServerError, err)
			return
		}

		jsonBytes, _ = json.Marshal(deps)

	} else {
		deps, err := models.GetAllDepartments()
		if err != nil {
			e.ResponseWithError(
				w, r, http.StatusInternalServerError, err)
			return
		}

		jsonBytes, _ = json.Marshal(deps)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
