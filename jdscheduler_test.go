package jdscheduler

import (
	"reflect"
	"testing"
	"time"
)

const timeLayout = "2006-Jan-02"

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Errorf("received: %v (type %v), expected: %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

///////////// SEASON TESTS //////////////

func TestNthSundayOfMonth(t *testing.T) {

	firstSunday, err := time.Parse(timeLayout, "2019-Sep-01")
	secondSunday, err := time.Parse(timeLayout, "2020-Aug-09")
	thirdSunday, err := time.Parse(timeLayout, "2020-Aug-16")
	fourthSunday, err := time.Parse(timeLayout, "2020-May-24")
	fifthSunday, err := time.Parse(timeLayout, "2020-May-31")

	date, err := nthSundayOfMonth(time.Month(9), 1, 2019)
	AssertEqual(t, firstSunday, date)

	date, err = nthSundayOfMonth(time.Month(8), 2, 2020)
	AssertEqual(t, secondSunday, date)

	date, err = nthSundayOfMonth(time.Month(8), 3, 2020)
	AssertEqual(t, thirdSunday, date)

	date, err = nthSundayOfMonth(time.Month(5), 4, 2020)
	AssertEqual(t, fourthSunday, date)

	date, err = nthSundayOfMonth(time.Month(5), 5, 2020)
	AssertEqual(t, fifthSunday, date)

	_, err = nthSundayOfMonth(time.Month(9), 6, 2019)
	if err == nil {
		t.Errorf("should have raised error")
	}

}

func TestSetOpenCloseWeeks(t *testing.T) {

	// should be sunday
	openPass, _ := time.Parse(timeLayout, "2020-Apr-19")
	closePass, _ := time.Parse(timeLayout, "2020-Oct-18")

	open, close, err := setOpenCloseWeeks(4, 3, 27, 2020, 7)
	AssertEqual(t, open, openPass)
	AssertEqual(t, close, closePass)

	open, close, err = setOpenCloseWeeks(4, 10, 6, 2020, 7)
	if err == nil {
		t.Errorf("should have raised nth sunday error")
	}
	open, close, err = setOpenCloseWeeks(4, 5, 7, 2020, 10)
	if err == nil {
		t.Errorf("should have raised too short error")
	}
}

func TestBlockTypeAssignment(t *testing.T) {

	open, _ := time.Parse(timeLayout, "2020-Apr-26")
	close, _ := time.Parse(timeLayout, "2020-Oct-04")

	opening1, _ := time.Parse(timeLayout, "2020-Apr-26")
	prime1, _ := time.Parse(timeLayout, "2020-Jun-07")
	prime12, _ := time.Parse(timeLayout, "2020-Aug-23")
	closing1, _ := time.Parse(timeLayout, "2020-Aug-30")

	AssertEqual(t, assignSeasonBlockType(opening1, open, close, 6), Opening)
	AssertEqual(t, assignSeasonBlockType(prime1, open, close, 6), Prime)
	AssertEqual(t, assignSeasonBlockType(prime12, open, close, 6), Prime)
	AssertEqual(t, assignSeasonBlockType(closing1, open, close, 6), Closing)

}

func TestSeasonInit(t *testing.T) {

	// error season
	_, errS := newSeason(time.July, 1, 3, 2020, 6)
	if errS == nil {
		t.Errorf("should have raised error: less than min weeks")
	}
}

func TestSeasonScheduling(t *testing.T) {

	var nextSeason func(sM time.Month, sW, nWk, year, nP int, sch *scheduler) (*season, error)
	nextSeason = func(sM time.Month, sW, nWk, year, nP int, sch *scheduler) (*season, error) {
		season, err := newSeason(sM, sW, nWk, year, nP)
		if err != nil {
			return nil, err
		}
		sch.assignSeason(season)
		return season, nil
	}

	// 6 ppl
	participantsSix := []string{"A", "B", "C", "D", "E", "F"}
	// long prime, double and single orders used in prime
	startM := time.April
	startW := 2
	nWk := 25
	schSixA := newScheduler(participantsSix)
	season, err := nextSeason(startM, startW, nWk, 2020, 6, schSixA)
	singleParticipantUsed := season.blocks[1].weeks[6*2].participant
	singleParticipantUsedIdx := indexOf(singleParticipantUsed, participantsSix)
	season, err = nextSeason(startM, startW, nWk, 2021, 6, schSixA)
	AssertEqual(t, participantsSix[(singleParticipantUsedIdx+1)%6], season.blocks[1].weeks[6*2].participant)

	// shorter prime, double and remaining orders used
	startM = time.April
	startW = 4
	nWk = 22
	schSixB := newScheduler(participantsSix)
	season, err = nextSeason(startM, startW, nWk, 2020, 6, schSixB)
	AssertEqual(t, season.blocks[1].weeks[7].participant, "D") // d was doubled
	AssertEqual(t, season.blocks[1].weeks[8].participant, "E") // everyone gets weeks in prime (remaining was used)
	AssertEqual(t, season.blocks[1].weeks[9].participant, "F")
	season, err = nextSeason(startM, startW, nWk, 2021, 6, schSixB)
	AssertEqual(t, season.blocks[1].weeks[0].participant, "E") // b2b should start where it left off
	AssertEqual(t, season.blocks[1].weeks[4].participant, "B") // should rotate when complete

	// should error
	startM = time.April
	startW = 4
	nWk = 4
	season, err = nextSeason(startM, startW, nWk, 2020, 6, schSixB)
	if err == nil {
		t.Errorf("should have raised error: min season error")
	}

}

///////////// BLOCK TESTS //////////////////
func TestBlockSegmentBlockWeeks(t *testing.T) {

	var checkSunday func(w week)
	checkSunday = func(w week) {
		if w.startDate.Weekday() != 0 {
			t.Errorf("found non Sunday start days: %s", w.startDate)
		}
	}

	sunday, err := time.Parse(timeLayout, "2020-Aug-02")
	if err != nil {
		t.Errorf("could not parse sample time")
	}

	bA := newBlock(sunday, sunday.AddDate(0, 0, 14), None) // 2 week
	bB := newBlock(sunday, sunday.AddDate(0, 0, 13), None) // 1.85 weeks -> should have 1 weeks
	bC := newBlock(sunday, sunday.AddDate(0, 0, 29), None) // 4.14 weeks -> should have 4 weeks

	// test the fallback
	if len(bA.weeks) != 2 || len(bB.weeks) != 1 || len(bC.weeks) != 4 {
		t.Errorf("\nincorrect week rounding\na:%d should be 2\nb:%d should be 1\nc:%d should be 4",
			len(bA.weeks), len(bB.weeks), len(bC.weeks))
	}

	// start days should always be sundays
	for _, w := range bA.weeks {
		if w.startDate.Weekday() != 0 {
			t.Errorf("found non Sunday start days: %s", w.startDate)
		}
	}
	for _, w := range bB.weeks {
		if w.startDate.Weekday() != 0 {
			t.Errorf("found non Sunday start days: %s", w.startDate)
		}
	}
	for _, w := range bC.weeks {
		checkSunday(w)
	}
}

func indexOf(str string, s []string) int {
	for i, v := range s {
		if str == v {
			return i
		}
	}
	return -1
}
