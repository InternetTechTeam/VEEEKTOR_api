package models

import (
	"VEEEKTOR_api/pkg/database/pgsql"
	e "VEEEKTOR_api/pkg/errors"
	"database/sql"
	"log"
	"time"
)

type NestedLab struct {
	Id           int       `json:"id"`
	CourseId     int       `json:"course_id"`
	Opens        time.Time `json:"opens"`  // UTC
	Closes       time.Time `json:"closes"` // UTC
	Topic        string    `json:"topic"`
	Requirements string    `json:"requirements,omitempty"`
	Example      string    `json:"example,omitempty"`
	LocationId   int       `json:"location_id,omitempty"`
	Attempts     int       `json:"attempts,omitempty"`
}

// Errors: ErrLabNotFound
func GetNestedLabById(labId int) (NestedLab, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id, opens, closes, topic, 
		requirements, example, location_id, attempts 
		FROM nested_labs WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	var lab NestedLab
	if err = stmt.QueryRow(&labId).Scan(
		&lab.Id, &lab.CourseId, &lab.Opens,
		&lab.Closes, &lab.Topic, &lab.Requirements,
		&lab.Example, &lab.LocationId, &lab.Attempts); err != nil {
		if err == sql.ErrNoRows {
			return lab, e.ErrNestedLabNotFound
		}
		log.Fatal(err)
	}

	return lab, nil
}

// Errors: ErrNestedLabsNotFound
func GetNestedLabsByCourseId(courseId int) ([]NestedLab, error) {
	stmt, err := pgsql.DB.Prepare(
		`SELECT id, course_id, opens, closes, topic 
		FROM nested_labs WHERE course_id = $1`)
	if err != nil {
		log.Fatal(err)
	}

	var labs []NestedLab
	var rows *sql.Rows
	if rows, err = stmt.Query(&courseId); err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var lab NestedLab
		if err = rows.Scan(
			&lab.Id, &lab.CourseId, &lab.Opens,
			&lab.Closes, &lab.Topic); err != nil {
			log.Fatal(err)
		}
		labs = append(labs, lab)
	}

	if len(labs) == 0 {
		return labs, e.ErrNestedLabsNotFound
	}

	return labs, nil
}

// Errors: ErrMissingFields, ErrNestedLabNotFound,
// ErrLocationNotFound, ErrCourseNotFound
func (lab *NestedLab) Validate() error {
	if len(lab.Topic) == 0 ||
		lab.Attempts == 0 ||
		lab.LocationId == 0 {
		return e.ErrMissingFields
	}

	if lab.Id != 0 {
		err := pgsql.DB.QueryRow(
			`SELECT 1 FROM nested_labs WHERE id = $1`,
			&lab.Id).Scan()
		if err != nil {
			if err == sql.ErrNoRows {
				return e.ErrNestedLabNotFound
			}
			log.Fatal(err)
		}
	}

	err := pgsql.DB.QueryRow(
		`SELECT 1 FROM locations WHERE id = $1`,
		&lab.LocationId).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return e.ErrLocationNotFound
		}
		log.Fatal(err)
	}

	err = pgsql.DB.QueryRow(
		`SELECT 1 FROM courses WHERE id = $1`,
		&lab.CourseId).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return e.ErrCourseNotFound
		}
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrMissingFields, ErrNestedLabNotFound,
// ErrLocationNotFound, ErrCourseNotFound
func (lab *NestedLab) Insert() error {
	if err := lab.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`INSERT INTO nested_labs(
		course_id, opens, closes, 
		topic, requirements, example, 
		location_id, attempts) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	_, err = stmt.Exec(&lab.CourseId, &lab.Opens, &lab.Closes,
		&lab.Topic, &lab.Requirements, &lab.Example,
		&lab.LocationId, &lab.Attempts)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Errors: ErrMissingFields, ErrNestedLabNotFound,
// ErrLocationNotFound, ErrCourseNotFound
func (lab *NestedLab) Update() error {
	if err := lab.Validate(); err != nil {
		return err
	}

	stmt, err := pgsql.DB.Prepare(
		`UPDATE nested_labs SET 
		course_id = $2, opens = $3, closes = $4, 
		topic = $5, requirements = $6, example = $7, 
		location_id = $8, attempts = $9
		WHERE id = $1`)
	if err != nil {
		log.Fatal(e.ErrCantPrepareDbStmt)
	}

	_, err = stmt.Exec(
		&lab.Id, &lab.CourseId, &lab.Opens,
		&lab.Closes, &lab.Topic, &lab.Requirements,
		&lab.Example, &lab.LocationId, &lab.Attempts)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func DeleteNestedLabById(labId int) error {
	stmt, err := pgsql.DB.Prepare(
		`DELETE FROM nested_labs WHERE id = $1`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(&labId)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
