package jdscheduler

/*
Scheduler - a scheduler for a set of participants that runs over many seasons
*/
type Scheduler struct {
	NParticipants int
	FairMap       map[string]int      // keep track of number of weeks for each participant
	pickOrder     map[BlockType]order // the order that participants pick a season block. holds across seasons
	pickIndex     map[BlockType]pick  // the index of the current picker in each block's pick order
}

type pick struct {
	single int // index of the current picker in pickOrder for blocks that only assign single weeks
	double int // index of the current picker in pickOrder for blocks that can assign back to back weeks
}

type order struct {
	single []string // order for assigning single weeks
	double []string // order for assigning double weeks
}

/*
NewScheduler - creats a new scheduler with a participants pick order for each blockType
*/
func NewScheduler(participants []string) *Scheduler {

	fm := make(map[string]int)
	for _, p := range participants {
		fm[p] = 0
	}

	// set intial pick order for each blockType. rotate starter for each block
	po := map[BlockType]order{
		Opening: order{
			rotatePickOrder(participants, 0),
			rotatePickOrder(participants, 0),
		},
		Prime: order{
			rotatePickOrder(participants, 0),
			rotatePickOrder(participants, 0),
		},
		Closing: order{
			rotatePickOrder(participants, 0),
			rotatePickOrder(participants, 0),
		},
	}

	pi := map[BlockType]pick{
		Opening: {0, 0},
		Prime:   {0, 0},
		Closing: {0, 0},
	}

	return &Scheduler{len(participants), fm, po, pi}
}

/*
AssignSeason assign season Block Weeks with a scheduler's state order of participants.
Each season block has its own pick order.
*/
func (sch *Scheduler) AssignSeason(s *Season) {
	for _, b := range s.Blocks {
		sch.assignBlockWeeks(&b)
	}
}

/*
Assigns a block's weeks with a scheduler's state order.
...
*/
func (sch *Scheduler) assignBlockWeeks(blk *Block) {

	blkType := blk.GetBlockType()
	weeks := blk.GetWeeks()
	pickIndex := sch.pickIndex[blkType]
	pickOrder := sch.pickOrder[blkType]

	nWpP := float32(len(weeks)) / float32(sch.NParticipants)
	// everyone can fit and more.
	if nWpP > 1 {

		// a participant gets 1 double week a block at most
		nDoubleWeeks := sch.NParticipants
		if nWpP < 2 {
			nDoubleWeeks = len(weeks) % sch.NParticipants
		}

		remaining := make([]string, sch.NParticipants)
		copy(remaining, rotatePickOrder(pickOrder.double, pickIndex.double))

		// assign all double weeks
		currWeek := 0
		for i := 0; i < nDoubleWeeks; i++ {
			participant := pickOrder.double[pickIndex.double]
			weeks[currWeek].AssignParticipant(participant)
			weeks[currWeek+1].AssignParticipant(participant)
			sch.FairMap[participant] += 2
			currWeek = (i + 1) * 2
			remaining = remove(remaining, participant)
			pickIndex.double++
			if pickIndex.double == sch.NParticipants {
				pickOrder.double = rotatePickOrder(pickOrder.double, 1)
				pickIndex.double = 0
			}
		}

		// some participants still need weeks in this block, use remaining
		if nDoubleWeeks < sch.NParticipants {
			for currWeek < len(weeks) && len(remaining) > 0 {
				participant := pop(remaining)
				weeks[currWeek].AssignParticipant(participant)
				sch.FairMap[participant]++
				currWeek++
			}
		} else {
			// just use single pick index
			for currWeek < len(weeks) {
				participant := pickOrder.single[pickIndex.single]
				weeks[currWeek].AssignParticipant(participant)
				sch.FairMap[participant]++
				pickIndex.single++
				currWeek++
				if pickIndex.single == sch.NParticipants {
					pickOrder.single = rotatePickOrder(pickOrder.single, 1)
					pickIndex.single = 0
				}
			}
		}

	} else {
		// participant gets less than or exactly 1 week. rotate across years and only assign single weeks
		for i := range weeks {
			participant := pickOrder.single[pickIndex.single]
			weeks[i].AssignParticipant(participant)
			sch.FairMap[participant]++
			pickIndex.single++
			// rotate and reset index when we have made a full rotation
			if pickIndex.single == sch.NParticipants {
				pickOrder.single = rotatePickOrder(pickOrder.single, 1)
				pickIndex.single = 0
			}
		}
	}

	sch.pickOrder[blkType] = pickOrder
	sch.pickIndex[blkType] = pickIndex
}

/* returns order that has been rotated by n steps as new slice */
func rotatePickOrder(order []string, steps int) []string {
	ret := make([]string, len(order))
	copy(ret, order)
	if len(ret) <= 0 || steps == 0 {
		return ret
	}
	steps = steps % len(ret)
	return append(ret[steps:], ret[:steps]...)
}

/* removes 1 string r from s */
func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

/*returns first element while removing from underlying s*/
func pop(s []string) string {
	pop := s[0]
	s = remove(s, pop)
	return pop
}

func minInt(vars ...int) (m int) {
	min := vars[0]
	for _, v := range vars {
		if v < min {
			min = v
		}
	}
	return min
}