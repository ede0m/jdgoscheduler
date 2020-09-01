package jdscheduler

import "time"

/*
Week represents 7 days in the schedule
*/
type Week struct {
	StartDate   time.Time
	Participant string
	PointValue  float32
}

/*
AssignParticipant assigns a participant p to the receiver week
*/
func (w *Week) AssignParticipant(p string) {
	w.Participant = p
}
