package ipefmt

import "regexp"

const (
	ArgColorNever  = "never"
	ArgColorAlways = "always"
	ArgColorAuto   = "auto"

	ArgTimeAcc = "accessed"
	ArgTimeMod = "modified"
	ArgTimeCrt = "created"

	ArgSortNone     = "none"
	ArgSortInode    = "inode"
	ArgSortMode     = "mode"
	ArgSortSize     = "size"
	ArgSortAccessed = "accessed"
	ArgSortModified = "modified"
	ArgSortCreated  = "created"
	ArgSortUser     = "user"
	ArgSortName     = "name"
)

type ArgsInfo struct {
	Across    bool
	All       bool
	Color     string
	Classify  bool
	Depth     uint8
	Filter    []*regexp.Regexp
	Header    bool
	Ignore    []*regexp.Regexp
	Inode     bool
	Long      bool
	OneLine   bool
	Reverse   bool
	Recursive bool
	Separator string
	Sort      string
	Sources   []string
	Time      []string
	Tree      bool
	Width     int
}
