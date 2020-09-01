package jdgoscheduler

import "time"

/*
Week represents 7 days in the schedule
*/
type week struct {
	startDate   time.Time
	participant string
	pointValue  float32
}
