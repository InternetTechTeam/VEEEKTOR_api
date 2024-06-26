package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"errors"
	"log"
)

type Department struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	EnvId int    `json:"env_id"`
}

// Errors: ErrDepsNotFound
func GetAllDepartments() ([]Department, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, name, env_id FROM departments`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var deps []Department
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var dep Department
		err = rows.Scan(&dep.Id, &dep.Name, &dep.EnvId)
		if err != nil {
			log.Fatal(err)
		}
		deps = append(deps, dep)
	}

	if len(deps) == 0 {
		return deps, e.ErrDepsNotFound
	}

	return deps, nil
}

// Errors: ErrDepNotFound
func GetDepartmentById(depId int) (Department, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, name, env_id FROM departments WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var dep Department
	err = stmt.QueryRow(&depId).Scan(&dep.Id, &dep.Name, &dep.EnvId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dep, e.ErrDepNotFound
		}
		log.Fatal(err)
	}

	return dep, nil
}

// Errors: ErrDepsNotFound
func GetAllDepartmentsByEnvironmentId(envId int) ([]Department, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, name, env_id FROM departments WHERE env_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var deps []Department
	rows, err := stmt.Query(&envId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var dep Department
		err = rows.Scan(&dep.Id, &dep.Name, &dep.EnvId)
		if err != nil {
			log.Fatal(err)
		}
		deps = append(deps, dep)
	}

	if len(deps) == 0 {
		return deps, e.ErrDepsNotFound
	}

	return deps, nil
}
