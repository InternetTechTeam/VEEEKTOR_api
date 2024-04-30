package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"errors"
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

// Errors: ErrCoursesNotFound
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
		if errors.Is(err, sql.ErrNoRows) {
			return course, e.ErrCourseNotFound
		}
		log.Fatal(err)
	}

	return course, nil
}

// Errors: ErrCoursesNotFound
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
		JOIN users AS u ON c.teacher_id=u.id 
		JOIN departments AS d_u ON d_u.id=u.dep_id 
		JOIN departments AS d_c ON d_c.id=c.dep_id
		WHERE c.id=$1`)
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
			if errors.Is(err, sql.ErrNoRows) {
				return courses, e.ErrCoursesNotFound
			}
			log.Fatal(err)
		}
		courses = append(courses, c)
	}

	if len(courses) == 0 {
		return courses, e.ErrCoursesNotFound
	}

	return courses, nil
}

// Errors: ErrCourseNotFound, ErrTermNotValid, ErrCourseNameNotValid
// ErrTeacherNotFound, ErrDepNotFound
func (c *Course) Validate() error {
	var exists bool
	if c.Id != 0 {
		err := pgsql.DB.QueryRow(
			`SELECT 1 FROM courses WHERE id=$1`,
			&c.Id).Scan(&exists)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return e.ErrCourseNotFound
			}
			log.Fatal(err)
		}
	}

	if c.Term <= 0 || c.Term > 14 {
		return e.ErrTermNotValid
	}

	if len(c.Name) == 0 {
		return e.ErrCourseNameNotValid
	}

	var roleId int
	err := pgsql.DB.QueryRow(
		`SELECT role_id from users WHERE id=$1`,
		&c.TeacherId).Scan(&roleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrTeacherNotFound
		}
		log.Fatal(err)
	}

	if roleId != 2 && roleId != 3 {
		return e.ErrTeacherNotFound
	}

	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM departments WHERE id=$1`,
		&c.DepId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.ErrDepNotFound
		}
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrCourseNotFound, ErrTermNotValid, ErrCourseNameNotValid
// ErrTeacherNotFound, ErrDepNotFound
func (c *Course) Insert() (int, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO courses 
		(name, term, teacher_id, markdown, dep_id) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if err = stmt.QueryRow(&c.Name, &c.Term,
		&c.TeacherId, &c.Markdown,
		&c.DepId).Scan(&c.Id); err != nil {
		log.Fatal(err)
	}

	if err = LinkUserWithCourse(c.TeacherId, c.Id); err != nil {
		log.Fatal(err)
	}

	return c.Id, nil
}

// Errors: ErrCourseNotFound, ErrTermNotValid, ErrCourseNameNotValid
// ErrTeacherNotFound, ErrDepNotFound
func (c *Course) Update() error {
	if c.Id == 0 {
		return e.ErrCourseIdNull
	}

	if err := c.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`UPDATE courses SET name=$2, term=$3, 
		teacher_id=$4, markdown=$5, dep_id=$6 
		WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err = stmt.Exec(
		&c.Id, &c.Name, &c.Term, &c.TeacherId,
		&c.Markdown, &c.DepId); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Expected that link will be unique.
// Errors: -
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

// Errors: -
func CheckUserBelongsToCourse(userId, courseId int) (bool, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT 1 FROM user_courses WHERE user_id=$1 AND course_id=$2`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var belongs bool
	if err = stmt.QueryRow(&userId, &courseId).Scan(&belongs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Fatal(err)
	}

	return true, nil
}
