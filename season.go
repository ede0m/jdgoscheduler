package jdgoscheduler

import (
	"errors"
	"time"
)

/*
Season for year is composed of blocks of weeks from open week to close week
*/
type season struct {
	year      int
	openWeek  time.Time
	closeWeek time.Time
	blocks    []block
}

/*
NewSeason creates a new season from
- a start month, week
- a close month, week
- number of participants
*/
func newSeason(openM time.Month, openW, nWeeks, year, nParticipants int) (s *season, err error) {
	open, close, err := setOpenCloseWeeks(openM, openW, nWeeks, year, nParticipants)
	if err != nil {
		return nil, err
	}
	weeks := initSeason(open, close, nParticipants)
	return &season{year, open, close, weeks}, nil
}

func initSeason(open, close time.Time, nParticipants int) []block {
	var blocks []block
	endItr := close.AddDate(0, 0, 7) // must surpass end date to include last block week
	currBt := Opening
	blkStart := open
	for d := open; d.Before(endItr) || d.Equal(endItr); d = d.AddDate(0, 0, 7) {
		bt := assignSeasonBlockType(d, open, close, nParticipants)
		if bt != currBt {
			blocks = append(blocks, newBlock(blkStart, d, currBt))
			blkStart = d
		}
		currBt = bt
	}
	return blocks
}

/*
	Sets the open and close weeks and check against participants (n).
	Each participant must get at least 1 week in the season.
	Weeks start on Sundays.
*/
func setOpenCloseWeeks(openM time.Month, openW, nWeeks, year, nP int) (o time.Time, c time.Time, err error) {

	startWeek, errE := nthSundayOfMonth(openM, openW, year)
	endWeek := startWeek.AddDate(0, 0, (7 * (nWeeks - 1)))
	if errE != nil {
		return o, c, errors.New("nth sunday error")
	}
	if nWeeks < nP {
		return o, c, errors.New("each participant must get at least 1 week in the season")
	}
	return startWeek, endWeek, nil
}

/*
	defines a week's block type within a season. returns None if out of bounds
*/
func assignSeasonBlockType(weekStart, open, close time.Time, nParticipants int) BlockType {

	weekYd := weekStart.YearDay()
	openYd := open.YearDay()
	closeYd := close.YearDay()

	endOpenYd := openYd + nParticipants*7
	startCloseYd := closeYd - nParticipants*7
	nWeeks := int(close.AddDate(0, 0, 7).Sub(open).Hours() / 24 / 7)

	// opening block will always have nParticipant weeks
	if openYd <= weekYd && weekYd < endOpenYd {
		return Opening
	}
	// we can set closing to nParticipants
	if nWeeks/nParticipants >= 2 {
		if startCloseYd < weekYd && weekYd <= closeYd {
			return Closing
		}
		// prime if not opening and closing
		return Prime
	}
	// if we can't fill closing with nParticipants, the rest will just be closing
	return Closing
}

/*
	Returns the date of the nth week in month of year
*/
func nthSundayOfMonth(month time.Month, wk, year int) (t time.Time, err error) {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	sundayDate := (8-int(firstDay.Weekday()))%7 + (7 * (wk - 1)) // 8 here to account for day 1 of month
	if sundayDate > 31 {
		return t, errors.New("no more than 5 weeks in a month")
	}
	return time.Date(year, month, sundayDate, 0, 0, 0, 0, time.UTC), nil
}
