package jdscheduler

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

const layoutISO = "2006-01-02"

/*
Schedule creates a new schedule full of seasons
*/
type Schedule struct {
	Years          int
	WeeksPerSeason int
	Participants   []string
	Seasons        []*season
	scheduler      scheduler
}

/*
NewSchedule - creats a new schedule
*/
func NewSchedule(startDate time.Time, nYears, weeksPerSeason int, participants []string) *Schedule {

	Seasons := make([]*season, nYears)
	schr := newScheduler(participants)
	startYear := startDate.Year()
	for y := startYear; y < startYear+nYears; y++ {
		season, err := newSeason(startDate, weeksPerSeason, len(participants))
		if err != nil {
			fmt.Println(err)
		}
		schr.assignSeason(season)
		Seasons[y-startYear] = season
	}

	return &Schedule{nYears, weeksPerSeason, participants, Seasons, *schr}
}

func (sch Schedule) String() string {
	var b bytes.Buffer
	var writer = tabwriter.NewWriter(&b, 0, 8, 0, '\t', tabwriter.AlignRight)

	for _, s := range sch.Seasons {
		fmt.Fprintln(writer, "open: ", s.openWeek.Format(layoutISO), "\tclose: ", s.closeWeek.Format(layoutISO), "\t", sch.WeeksPerSeason, "wk")
		fmt.Fprintln(writer)
		for _, b := range s.blocks {
			for _, w := range b.GetWeeks() {
				fmt.Fprintln(writer, w.startDate.Format(layoutISO), "\t", b.GetBlockType(), "\t", w.participant)
			}
		}
		fmt.Fprintln(writer)
		for k, v := range sch.scheduler.fairMap {
			fmt.Fprintf(writer, "participant [%s] weeks [%d]\n", k, v)
		}
		fmt.Fprintln(writer)
		fmt.Fprintln(writer, "-------------------------------------------------------")
		fmt.Fprintln(writer)
	}
	writer.Flush()
	return b.String()
}
