package jdscheduler

import "time"

/*
Block is a set of weeks within a season defined by some type
*/
type Block struct {
	open, close time.Time
	blockType   BlockType
	weeks       []Week
}

/*
BlockType defines the type of block within a season
*/
type BlockType int

/* Types of Block types within a season */
const (
	None BlockType = iota
	Opening
	Prime
	Closing
)

var blockTypeStrMap = map[BlockType]string{
	None:    "None",
	Opening: "Opening",
	Prime:   "Prime",
	Closing: "Closing",
}

/* BlockType implements Stringer */
func (b BlockType) String() string {
	t, exists := blockTypeStrMap[b]
	if !exists {
		return blockTypeStrMap[None]
	}
	return t
}

/*
NewBlock creates a new Block between two dates
*/
func NewBlock(blkStartDt, blkEndDt time.Time, bt BlockType) Block {
	return Block{blkStartDt, blkEndDt, bt, segmentBlockWeeks(blkStartDt, blkEndDt)}
}

/*
GetOpenClose gets the open and close weeks for a block
*/
func (b Block) GetOpenClose() (time.Time, time.Time) {
	return b.open, b.close
}

/*
GetBlockType gets the block type of a block
*/
func (b Block) GetBlockType() BlockType {
	return b.blockType
}

/*
GetWeeks gets the weeks of a block
*/
func (b Block) GetWeeks() []Week {
	return b.weeks
}

/*
	Creates Weeks within Block. Block can only have a whole number of weeks.
	It will fall back (floor) the number of weeks is float
*/
func segmentBlockWeeks(blkStartDt, blkEndDt time.Time) []Week {
	var weeks []Week
	for d := blkStartDt; d.Before(blkEndDt) || d.Equal(blkEndDt); d = d.AddDate(0, 0, 7) {
		// if we surpass end date, we fall back to whole number of weeks in block
		if !d.AddDate(0, 0, 7).After(blkEndDt) {
			weeks = append(weeks, Week{d, "", 0})
		}
	}
	return weeks
}
