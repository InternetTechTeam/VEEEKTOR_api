package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"log"
)

type NestedInfo struct {
	Id       int    `json:"id"`
	CourseId int    `json:"course_id"`
	Markdown string `json:"markdown,omitempty"`
}

// Errors: ErrNestedInfoNotFound
func GetNestedInfoById(infoId int) (NestedInfo, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id, markdown 
		FROM nested_info WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var info NestedInfo
	if err = stmt.QueryRow(&infoId).Scan(
		&info.Id, &info.CourseId, &info.Markdown); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return info, e.ErrNestedInfoNotFound
	}

	return info, nil
}

// Errors: ErrNestedInfosNotFound
func GetNestedInfosByCourseId(courseId int) ([]NestedInfo, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id 
		FROM nested_info WHERE course_id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var infos []NestedInfo
	var rows *sql.Rows
	if rows, err = stmt.Query(&courseId); err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var info NestedInfo
		if err = rows.Scan(
			&info.Id, &info.CourseId); err != nil {
			log.Fatal(err)
		}
		infos = append(infos, info)
	}

	if len(infos) == 0 {
		return infos, e.ErrNestedInfosNotFound
	}

	return infos, nil
}

// Errors: ErrCourseNotFound, ErrNestedInfoNotFound
func (info *NestedInfo) Validate() error {
	if info.Id != 0 {
		stmt, err := pgsql.DB.Prepare(
			`SELECT 1 FROM nested_info WHERE id=$1`)
		if err != nil {
			log.Fatal(err)
		}
		var exist bool
		if err = stmt.QueryRow(&info.Id).Scan(&exist); err != nil {
			if err == sql.ErrNoRows {
				return e.ErrNestedInfoNotFound
			}
			log.Fatal(err)
		}
	}

	if _, err := GetCourseById(info.CourseId); err != nil {
		return e.ErrCourseNotFound
	}
	return nil
}

// Errors: ErrCourseNotFound
func (info *NestedInfo) Insert() error {
	if err := info.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO nested_info(course_id, markdown)
		VALUES ($1, $2)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(info.CourseId, info.Markdown); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrCourseNotFound, ErrMissingFields, ErrNestedInfoNotFound
func (info *NestedInfo) Update() error {
	if info.Id == 0 {
		return e.ErrMissingFields
	}

	if err := info.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`UPDATE nested_info SET course_id = $2, markdown = $3
		WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(
		&info.Id, &info.CourseId, &info.Markdown); err != nil {
		log.Fatal(err)
	}

	return nil
}

func DeleteNestedInfoById(infoId int) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM nested_info WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err = stmt.Exec(&infoId); err != nil {
		log.Fatal(err)
	}

	return nil
}
