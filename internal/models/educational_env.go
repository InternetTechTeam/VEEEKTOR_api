package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"log"
)

type EducationalEnv struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllEducationalEnvs() ([]EducationalEnv, error) {
	// First educational environment supposed to be for admins
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, name from educational_envs WHERE id != 1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var envs []EducationalEnv
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var env EducationalEnv
		err = rows.Scan(&env.Id, &env.Name)
		if err != nil {
			log.Fatal(err)
		}
		envs = append(envs, env)
	}

	return envs, nil
}

func GetEducationalEnvironmentById(envId int) (EducationalEnv, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, name FROM educational_envs WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var env EducationalEnv
	err = stmt.QueryRow(&envId).Scan(&env.Id, &env.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return env, e.ErrEdEnvNotFound
		}
		log.Fatal(err)
	}

	return env, nil
}
