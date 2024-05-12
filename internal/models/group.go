package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"errors"
	"log"
)

type Group struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	DepId int    `json:"dep_id"`
}

func GetGroupById(groupId int) (Group, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, name, dep_id FROM groups WHERE id=$1`)
	if err != nil {
		log.Fatal(err)
	}

	var g Group
	if err = stmt.QueryRow(&groupId).Scan(
		&g.Id, &g.Name, &g.DepId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return g, e.ErrGroupNotFound
		}
		log.Fatal(err)
	}

	return g, nil
}

// Errors: ErrGroupsNotFound
func GetAllGroupsByDepId(depId int) ([]Group, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, name, dep_id FROM groups WHERE dep_id=$1`)
	if err != nil {
		log.Fatal(err)
	}

	var rows *sql.Rows
	if rows, err = stmt.Query(&depId); err != nil {
		log.Fatal(err)
	}

	var groups []Group
	for rows.Next() {
		var g Group
		if err = rows.Scan(&g.Id, &g.Name, &g.DepId); err != nil {
			log.Fatal(err)
		}
		groups = append(groups, g)
	}

	if len(groups) == 0 {
		return groups, e.ErrGroupsNotFound
	}

	return groups, nil
}

// Errors: ErrMissingFields, ErrDepNotFound
func (g *Group) Validate() error {
	if g.Name == "" || g.DepId == 0 {
		return e.ErrMissingFields
	}

	var exists bool
	err := pgsql.DB.QueryRow(
		`SELECT 1 FROM departments WHERE id=$1`,
		&g.DepId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.ErrDepNotFound
		}
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrMissingFields, ErrDepNotFound
func (g *Group) Insert() error {
	if err := g.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO groups(name, dep_id) VALUES ($1, $2)`)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = stmt.Exec(&g.Name, &g.DepId); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Errors: -
func DeleteGroupById(groupId int) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM groups WHERE id=$1`)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = stmt.Exec(&groupId); err != nil {
		log.Fatal(err)
	}

	return nil
}
