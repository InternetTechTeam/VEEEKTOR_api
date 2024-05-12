package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"errors"
	"log"

	"github.com/golang-jwt/jwt"
)

type NestedInfo struct {
	Id       int    `json:"id"`
	CourseId int    `json:"course_id"`
	Name     string `json:"name"`
	Markdown string `json:"markdown,omitempty"`
}

// Errors: ErrNestedInfoNotFound
func GetNestedInfoById(infoId int) (NestedInfo, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id, name, markdown 
		FROM nested_infos WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var info NestedInfo
	if err = stmt.QueryRow(&infoId).Scan(&info.Id,
		&info.CourseId, &info.Name, &info.Markdown); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return info, e.ErrNestedInfoNotFound
		}
		log.Fatal(err)
	}

	return info, nil
}

// Errors: ErrNestedInfosNotFound
func GetNestedInfosByCourseId(courseId int) ([]NestedInfo, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id, name
		FROM nested_infos WHERE course_id=$1`)
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
			&info.Id, &info.CourseId, &info.Name); err != nil {
			log.Fatal(err)
		}
		infos = append(infos, info)
	}

	if len(infos) == 0 {
		return infos, e.ErrNestedInfosNotFound
	}

	return infos, nil
}

// Errors: ErrCourseNotFound, ErrNestedInfoNotFound, ErrMissingFields
func (info *NestedInfo) Validate() error {
	if len(info.Name) == 0 {
		return e.ErrMissingFields
	}
	var exists bool

	if info.Id != 0 {
		err := pgsql.DB.QueryRow(
			`SELECT 1 FROM nested_infos WHERE id=$1`,
			&info.Id).Scan(&exists)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return e.ErrNestedLabNotFound
			}
			log.Fatal(err)
		}
	}

	err := pgsql.DB.QueryRow(
		`SELECT 1 FROM courses WHERE id=$1`,
		&info.CourseId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrCourseNotFound
		}
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrCourseNotFound, ErrNestedInfoNotFound, ErrMissingFields
func (info *NestedInfo) Insert() error {
	if err := info.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO nested_infos(course_id, name, markdown)
		VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(&info.CourseId, &info.Name, &info.Markdown); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrCourseNotFound, ErrMissingFields,
// ErrNestedInfoNotFound, ErrMissingFields
func (info *NestedInfo) Update() error {
	if info.Id == 0 {
		return e.ErrMissingFields
	}

	if err := info.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`UPDATE nested_infos 
		SET course_id=$2, name=$3, markdown=$4
		WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err := stmt.Exec(
		&info.Id, &info.CourseId, &info.Name, &info.Markdown); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Errors: -
func DeleteNestedInfoById(infoId int) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM nested_infos WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err = stmt.Exec(&infoId); err != nil {
		log.Fatal(err)
	}

	return nil
}

// User have: 0 - no access, 1 - read access, 2 - write access
func (i *NestedInfo) CheckAccess(claims jwt.MapClaims) int {
	var teacherId int
	err := pgsql.DB.QueryRow(
		`SELECT teacher_id FROM 
		courses WHERE id=$1`, &i.CourseId).Scan(&teacherId)
	if err != nil {
		log.Fatal(err)
	}

	if claims["user_id"].(int) == teacherId {
		return 2
	}

	var exists int
	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM group_courses 
		WHERE group_id=$1 and course_id=$2`,
		claims["group_id"].(int), &i.CourseId).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	}

	return exists
}
