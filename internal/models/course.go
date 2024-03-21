package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"log"
)

type Course struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Term      int    `json:"term"`
	TeacherId int    `json:"teacher_id"`
	Markdown  string `json:"markdown,omitempty"`
	DepId     int    `json:"dep_id"`
}

func GetAllCoursesByUserId(userId int) ([]Course, error) {
	selCourseIdStmt, err := pgsql.DB.Prepare(
		`SELECT course_id from user_courses WHERE user_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	selCourseStmt, err := pgsql.DB.Prepare(
		`SELECT name, term, teacher_id, dep_id 
		 FROM courses WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var rows *sql.Rows
	if rows, err = selCourseIdStmt.Query(&userId); err != nil {
		log.Fatal(err)
	}
	var courses []Course
	for rows.Next() {
		var course Course

		if err = rows.Scan(&course.Id); err != nil {
			log.Fatal(err)
		}

		if err := selCourseStmt.QueryRow(&course.Id).Scan(
			&course.Name, &course.Term,
			&course.TeacherId, &course.DepId); err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
			}
			return courses, e.ErrCoursesNotFound
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func GetCourseById(courseId int) (Course, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT name, term, teacher_id, markdown, dep_id 
		FROM courses WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var course Course
	course.Id = courseId
	if err := stmt.QueryRow(&courseId).Scan(
		&course.Name, &course.Term, &course.TeacherId,
		&course.Markdown, &course.DepId); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return course, e.ErrCoursesNotFound
	}

	return course, nil
}
