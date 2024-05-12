package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
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
	ModifiedAt int64 `json:"modified_at"`
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
func GetAllCoursesByGroupId(groupId int) ([]CourseMultipleExportDTO, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT c.id, c.name, c.term, d_c.name, u.name, 
		u.patronymic, u.surname, d_u.name, c.modified_at 
		FROM courses AS c 
		JOIN users AS u ON c.teacher_id=u.id 
		JOIN departments AS d_u ON d_u.id=u.dep_id 
		JOIN departments AS d_c ON d_c.id=c.dep_id
		JOIN group_courses AS gc ON gc.group_id=$1
		WHERE c.id=gc.course_id`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var courses []CourseMultipleExportDTO
	var rows *sql.Rows
	if rows, err = stmt.Query(&groupId); err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var c CourseMultipleExportDTO
		var t time.Time
		if err := rows.Scan(
			&c.Id, &c.Name, &c.Term, &c.Dep, &c.Teacher.Name,
			&c.Teacher.Patronymic, &c.Teacher.Surname,
			&c.Teacher.Dep, &t); err != nil {
			log.Fatal(err)
		}
		c.ModifiedAt = t.Unix()
		courses = append(courses, c)
	}

	if len(courses) == 0 {
		return courses, e.ErrCoursesNotFound
	}

	return courses, nil
}

// Errors: ErrCoursesNotFound
func GetAllCoursesByTeacherId(teacherId int) ([]CourseMultipleExportDTO, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT c.id, c.name, c.term, d_c.name, u.name, 
		u.patronymic, u.surname, d_u.name, c.modified_at 
		FROM courses AS c 
		JOIN users AS u ON c.teacher_id=u.id 
		JOIN departments AS d_u ON d_u.id=u.dep_id 
		JOIN departments AS d_c ON d_c.id=c.dep_id
		WHERE c.teacher_id=$1`)
	if err != nil {
		log.Fatal(err)
	}

	var courses []CourseMultipleExportDTO
	var rows *sql.Rows
	if rows, err = stmt.Query(&teacherId); err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var c CourseMultipleExportDTO
		var t time.Time
		if err := rows.Scan(
			&c.Id, &c.Name, &c.Term, &c.Dep, &c.Teacher.Name,
			&c.Teacher.Patronymic, &c.Teacher.Surname,
			&c.Teacher.Dep, &t); err != nil {
			log.Fatal(err)
		}
		c.ModifiedAt = t.Unix()
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
		teacher_id=$4, markdown=$5, dep_id=$6, 
		modified_at=$7
		WHERE id=$1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	if _, err = stmt.Exec(
		&c.Id, &c.Name, &c.Term, &c.TeacherId,
		&c.Markdown, &c.DepId, time.Now()); err != nil {
		log.Fatal(err)
	}

	return nil
}

// User have: 0 - no access, 1 - read access, 2 - write access
func (c *Course) CheckAccess(claims jwt.MapClaims) int {
	if c.TeacherId == 0 {
		err := pgsql.DB.QueryRow(
			`SELECT teacher_id FROM courses WHERE id=$1`,
			&c.Id).Scan(&c.TeacherId)
		if err != nil {
			log.Fatal(err)
		}
	}

	if c.TeacherId == claims["user_id"].(int) {
		return 2
	}

	var exists int
	err := pgsql.DB.QueryRow(
		`SELECT 1 FROM group_courses WHERE group_id=$1 and course_id=$2`,
		claims["group_id"].(int), &c.Id).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatal(err)
	}

	return exists
}
