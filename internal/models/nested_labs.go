package models

import "time"

type NestedLab struct {
	Id           int
	CourseId     int
	Opens        time.Time
	Closes       time.Time
	Topic        string
	Requirements string
	Example      string
	LocationId   int
	Attempts     int
}
