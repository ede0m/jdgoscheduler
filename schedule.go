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
	Years          int       `json:"years"`
	UnitsPerSeason int       `json:"unitsPerSeason"`
	Participants   []string  `json:"participants"`
	Seasons        []*Season `json:"seasons"`
	scheduler      scheduler
}

/*
NewSchedule - creats a new schedule
*/
func NewSchedule(start time.Time, nYears, unitsPerSeason int, participants []string) (*Schedule, error) {

	Seasons := make([]*Season, nYears)
	schr := newScheduler(participants)
	startYear := start.Year()
	startMonth := start.Month()
	startDay := start.Day()
	for y := startYear; y < startYear+nYears; y++ {
		snStartDate := time.Date(y, startMonth, startDay, 0, 0, 0, 0, time.UTC)
		season, err := newSeason(snStartDate, unitsPerSeason, len(participants))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		schr.assignSeason(season)
		Seasons[y-startYear] = season
	}

	return &Schedule{nYears, unitsPerSeason, participants, Seasons, *schr}, nil
}

func (sch Schedule) String() string {
	var b bytes.Buffer
	var writer = tabwriter.NewWriter(&b, 0, 8, 0, '\t', tabwriter.AlignRight)

	for _, s := range sch.Seasons {
		fmt.Fprintln(writer, "open: ", s.OpenWeek.Format(layoutISO), "\tclose: ", s.CloseWeek.Format(layoutISO), "\t", sch.UnitsPerSeason, "wk")
		fmt.Fprintln(writer)
		for _, b := range s.Blocks {
			for _, u := range b.Units {
				fmt.Fprintln(writer, u.Start.Format(layoutISO), "\t", b.BlockType, "\t", u.Participant)
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
