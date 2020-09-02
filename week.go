package jdscheduler

import "time"

/*
Week represents 7 days in the schedule
*/
type Week struct {
	StartDate   time.Time
	Participant string
	pointValue  float32
}
