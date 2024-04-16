package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"log"
	"time"
)

type NestedTest struct {
	Id         int           `json:"id"`
	CourseId   int           `json:"course_id"`
	Opens      time.Time     `json:"opens"`
	Closes     time.Time     `json:"closes"`
	TasksCount int           `json:"tasks_count,omitempty"`
	Topic      string        `json:"topic"`
	LocationId int           `json:"location_id,omitempty"`
	Attempts   int           `json:"attempts,omitempty"`
	Password   string        `json:"password,omitempty"`
	TimeLimit  time.Duration `json:"time_limit,omitempty"`
}

// Errors: ErrTestNotFound
func GetNestedTestById(testId int) (NestedTest, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id, opens, closes, 
		tasks_count, topic, location_id, 
		attempts, password, time_limit 
		FROM nested_tests WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var duration string
	var test NestedTest
	if err = stmt.QueryRow(&testId).Scan(
		&test.Id, &test.CourseId, &test.Opens,
		&test.Closes, &test.TasksCount, &test.Topic,
		&test.LocationId, &test.Attempts, &test.Password,
		&duration); err != nil {
		if err == sql.ErrNoRows {
			return test, e.ErrNestedTestNotFound
		}
		log.Fatal(err)
	}

	return test, nil
}

// Errors: ErrNestedTestsNotFound
func GetNestedTestsByCourseId(courseId int) ([]NestedTest, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id, opens, closes, topic 
		FROM nested_tests WHERE course_id = $1`)
	if err != nil {
		log.Fatal(err)
	}

	var tests []NestedTest
	var rows *sql.Rows
	if rows, err = stmt.Query(&courseId); err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var test NestedTest
		if err = rows.Scan(
			&test.Id, &test.CourseId, &test.Opens,
			&test.Closes, &test.Topic); err != nil {
			log.Fatal(err)
		}
		tests = append(tests, test)
	}

	if len(tests) == 0 {
		return tests, e.ErrNestedTestsNotFound
	}

	return tests, nil
}

// Errors: ErrMissingFields, ErrTimeLimitTooShort,
// ErrNestedTestNotFound, ErrLocationNotFound, ErrCourseNotFound
func (test *NestedTest) Validate() error {
	if len(test.Topic) == 0 ||
		test.Attempts == 0 ||
		test.LocationId == 0 ||
		int(test.TimeLimit.Seconds()) == 0 {
		return e.ErrMissingFields
	}

	if int(test.TimeLimit.Minutes()) < 5 {
		return e.ErrTimeLimitTooShort
	}

	if test.Id != 0 {
		err := pgsql.DB.QueryRow(
			`SELECT 1 FROM nested_tests WHERE id = $1`,
			&test.Id).Scan()
		if err != nil {
			if err == sql.ErrNoRows {
				return e.ErrNestedTestNotFound
			}
			log.Fatal(err)
		}
	}

	err := pgsql.DB.QueryRow(
		`SELECT 1 FROM locations WHERE id = $1`,
		&test.LocationId).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return e.ErrLocationNotFound
		}
		log.Fatal(err)
	}

	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM courses WHERE id = $1`,
		&test.CourseId).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return e.ErrCourseNotFound
		}
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrMissingFields, ErrTimeLimitTooShort,
// ErrNestedTestNotFound, ErrLocationNotFound, ErrCourseNotFound
func (test *NestedTest) Insert() error {
	if err := test.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO nested_tests(
		course_id, opens, closes, 
		tasks_count, topic, location_id, 
		attempts, password, time_limit) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	_, err = stmt.Exec(
		&test.CourseId, &test.Opens, &test.Closes,
		&test.TasksCount, &test.Topic, &test.LocationId,
		&test.Attempts, &test.Password, &test.TimeLimit)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrMissingFields, ErrTimeLimitTooShort,
// ErrNestedTestNotFound, ErrLocationNotFound, ErrCourseNotFound
func (test *NestedTest) Update() error {
	if err := test.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`UPDATE nested_tests SET 
		course_id = $2, opens = $3, closes = $4, 
		tasks_count = $5, topic = $6, location_id = $7, 
		attempts = $8, password = $9, time_limit = $10
		WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	_, err = stmt.Exec(
		&test.Id, &test.CourseId, &test.Opens,
		&test.Closes, &test.TasksCount, &test.Topic,
		&test.LocationId, &test.Attempts, &test.Password,
		&test.TimeLimit)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func DeleteNestedTestById(testId int) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM nested_tests WHERE id = $1`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(&testId)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
