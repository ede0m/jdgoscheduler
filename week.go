package jdscheduler

import (
	"time"

	"github.com/google/uuid"
)

/*
Week represents 7 days in the schedule
*/
type Week struct {
	ID          uuid.UUID `json:"id"`
	StartDate   time.Time `json:"startDate"`
	Participant string    `json:"participant"`
	pointValue  float32   // TODO: point scheme eventually?
}
