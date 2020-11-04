package jdscheduler

import (
	"time"

	"github.com/google/uuid"
)

/*
Week represents 7 days in the schedule
*/
type Unit struct {
	ID          uuid.UUID `json:"id"`
	Start       time.Time `json:"start"`
	Participant string    `json:"participant"`
	pointValue  float32   // TODO: point scheme eventually?
}
