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

type CourseMultipleExportDTO struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Term    int    `json:"term"`
	Dep     string `json:"department"`
	Teacher struct {
		Name       string `json:"name"`
		Patronymic string `json:"patronymic"`
		Surname    string `json:"surname"`
		Dep        string `json:"department"`
	} `json:"teacher"`
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

func GetAllCoursesByUserId(userId int) ([]CourseMultipleExportDTO, error) {
	selCourseIdStmt, err := pgsql.DB.Prepare(
		`SELECT course_id from user_courses WHERE user_id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	selCourseStmt, err := pgsql.DB.Prepare(
		`SELECT c.name, c.term, d_c.name, 
		u.name, u.patronymic, u.surname, d_u.name 
		FROM courses AS c 
		JOIN users AS u ON c.teacher_id = u.id 
		JOIN departments AS d_u ON d_u.id = u.dep_id 
		JOIN departments AS d_c ON d_c.id = c.dep_id
		WHERE c.id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var rows *sql.Rows
	if rows, err = selCourseIdStmt.Query(&userId); err != nil {
		log.Fatal(err)
	}
	var courses []CourseMultipleExportDTO
	for rows.Next() {
		var c CourseMultipleExportDTO

		if err = rows.Scan(&c.Id); err != nil {
			log.Fatal(err)
		}

		if err := selCourseStmt.QueryRow(&c.Id).Scan(
			&c.Name, &c.Term, &c.Dep, &c.Teacher.Name,
			&c.Teacher.Patronymic, &c.Teacher.Surname, &c.Teacher.Dep); err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
			}
			return courses, e.ErrCoursesNotFound
		}
		courses = append(courses, c)
	}

	return courses, nil
}

func (c *Course) Validate() error {
	if c.Term <= 0 || c.Term > 14 {
		return e.ErrTermNotValid
	}

	if len(c.Name) == 0 {
		return e.ErrCourseNameInvalid
	}

	stmt, err := pgsql.DB.Prepare(
		`SELECT role_id from users WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var roleId int
	if stmt.QueryRow(&c.TeacherId).Scan(&roleId); roleId != 2 && roleId != 3 {
		return e.ErrTeacherNotFound
	}

	if _, err = GetDepartmentById(c.DepId); err != nil {
		return err
	}

	return nil
}

func (c *Course) Insert() error {
	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO courses 
		(name, term, teacher_id, markdown, dep_id) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if err := c.Validate(); err != nil {
		return err
	}

	if err = stmt.QueryRow(&c.Name, &c.Term,
		&c.TeacherId, &c.Markdown,
		&c.DepId).Scan(&c.Id); err != nil {
		log.Fatal(err)
	}

	if err = LinkUserWithCourse(c.TeacherId, c.Id); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Expected that link will be unique. Any error cause Fatal.
func LinkUserWithCourse(userId, courseId int) error {
	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO user_courses (user_id, course_id) VALUES ($1, $2)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err = stmt.Exec(userId, courseId); err != nil {
		log.Fatal(err)
	}

	return nil
}
