package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"log"
)

type Course struct {
	Id        int    `josn:"id"`
	Name      string `josn:"name"`
	TeacherId int    `josn:"teacher_id"`
	Markdown  string `josn:"markdown"`
	DepId     int    `josn:"dep_id"`
}

func GetCoursesByUserId(userId int) ([]Course, error) {
	selCourseIdStmt, err := pgsql.DB.Prepare(
		`SELECT course_id from user_courses WHERE user_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	selCourseStmt, err := pgsql.DB.Prepare(
		`SELECT name, teacher_id, markdown, dep_id 
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
			&course.Name, &course.TeacherId,
			&course.Markdown, &course.DepId); err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
			}
			return courses, e.ErrCoursesNotFound
		}
		courses = append(courses, course)
	}

	return courses, nil
}
