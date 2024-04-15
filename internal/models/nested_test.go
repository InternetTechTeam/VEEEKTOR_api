package models

import "time"

type NestedTest struct {
	Id         int
	CourseId   int
	Opens      time.Time
	Closes     time.Time
	TasksCount int
	Topic      string
	LocationId int
	Attempts   int
	Password   string
	TimeLimit  time.Time
}
