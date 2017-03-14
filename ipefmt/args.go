package ipefmt

import "regexp"

const (
	ArgColorNever  = "never"
	ArgColorAlways = "always"
	ArgColorAuto   = "auto"

	ArgTimeAcc = "accessed"
	ArgTimeMod = "modified"
	ArgTimeCrt = "created"
)

type ArgsInfo struct {
	Across    bool
	All       bool
	Color     string
	Classify  bool
	Depth     uint8
	Filter    *regexp.Regexp
	Header    bool
	Ignore    *regexp.Regexp
	Inode     bool
	Long      bool
	OneLine   bool
	Reverse   bool
	Recursive bool
	Separator string
	Sources   []string
	Time      []string
	Tree      bool
	Width     int
}
