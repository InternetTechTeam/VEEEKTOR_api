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

type GroupCourse struct {
	Id       int `json:"id"`
	GroupId  int `json:"group_id"`
	CourseId int `json:"course_id"`
}

// Errors: ErrMissingFields, ErrGroupNotExist, ErrCourseNotFound
func (gc *GroupCourse) Validate() error {
	if gc.GroupId == 0 || gc.CourseId == 0 {
		return e.ErrMissingFields
	}

	var exists bool
	err := pgsql.DB.QueryRow(
		`SELECT 1 FROM groups WHERE id=$1`,
		&gc.GroupId).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	}
	if !exists {
		return e.ErrGroupNotExist
	}

	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM courses WHERE id=$1`,
		&gc.CourseId).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	}
	if !exists {
		return e.ErrCourseNotFound
	}

	return nil
}

// Errors: ErrGroupLinkedToCourse, ErrMissingFields,
// ErrGroupNotExist, ErrCourseNotFound
func (gc *GroupCourse) Insert() error {
	var exists bool
	err := pgsql.DB.QueryRow(
		`SELECT 1 from group_courses WHERE 
		group_id=$1 AND course_id=$2`,
		&gc.GroupId, &gc.CourseId).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	}
	if exists {
		return e.ErrGroupLinkedToCourse
	}

	if err = gc.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO group_courses(group_id, course_id)
		VALUES ($1, $2)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(&gc.GroupId, &gc.CourseId); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (gc *GroupCourse) Delete() error {
	var id int
	err := pgsql.DB.QueryRow(
		`SELECT id from group_courses WHERE 
		group_id=$1 AND course_id=$2`,
		&gc.GroupId, &gc.CourseId).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	}
	if id == 0 {
		return e.ErrGroupNotLinkedToCourse
	}

	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM group_courses WHERE id=$1`)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = stmt.Exec(&id); err != nil {
		log.Fatal(err)
	}

	return nil
}
