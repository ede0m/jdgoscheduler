package jdscheduler

import "time"

/*
Week represents 7 days in the schedule
*/
type Week struct {
	StartDate   time.Time `json:"startDate"`
	Participant string    `json:"participant"`
	pointValue  float32   // TODO: point scheme eventually?
}
