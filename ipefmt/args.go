package ipefmt

import "regexp"

const (
	// ArgColorNever represents an option for the `color` flag.
	// It means the output will never be printed with colors.
	ArgColorNever = "never"
	// ArgColorAlways represents an option for the `color` flag.
	// It means the output will always be printed with colors.
	ArgColorAlways = "always"
	// ArgColorAuto represents an option for the `color` flag.
	// It means the output will be printed with colors, only if it is stdout.
	ArgColorAuto = "auto"

	// ArgTimeAcc represents an option for the `time` flag.
	// It means that the "accessed time" will be printed in long view.
	ArgTimeAcc = "accessed"
	// ArgTimeMod represents an option for the `time` flag.
	// It means that the "modified time" will be printed in long view.
	ArgTimeMod = "modified"
	// ArgTimeCrt represents an option for the `time` flag.
	// It means that the "created time" will be printed in long view.
	ArgTimeCrt = "created"

	// ArgSortNone represents an option for the `sort` flag.
	// It means the output will not be sorted.
	ArgSortNone = "none"
	// ArgSortInode represents an option for the `sort` flag.
	// It means the output will be sorted by inode.
	ArgSortInode = "inode"
	// ArgSortMode represents an option for the `sort` flag.
	// It means the output will be sorted by mode.
	ArgSortMode = "mode"
	// ArgSortSize represents an option for the `sort` flag.
	// It means the output will be sorted by size.
	ArgSortSize = "size"
	// ArgSortLinks represents an option for the `sort` flag.
	// It means the output will be sorted by link.
	ArgSortLinks = "link"
	// ArgSortBlocks represents an option for the `sort` flag.
	// It means the output will be sorted by blocks.
	ArgSortBlocks = "blocks"
	// ArgSortAccessed represents an option for the `sort` flag.
	// It means the output will be sorted by accessed time.
	ArgSortAccessed = "accessed"
	// ArgSortModified represents an option for the `sort` flag.
	// It means the output will be sorted by modified time.
	ArgSortModified = "modified"
	// ArgSortCreated represents an option for the `sort` flag.
	// It means the output will be sorted by created time.
	ArgSortCreated = "created"
	// ArgSortUser represents an option for the `sort` flag.
	// It means the output will be sorted by user.
	ArgSortUser = "user"
	// ArgSortGroup represents an option for the `sort` flag.
	// It means the output will be sorted by group.
	ArgSortGroup = "group"
	// ArgSortName represents an option for the `sort` flag.
	// It means the output will be sorted by name.
	ArgSortName = "name"
)

// ArgsInfo represents all the arguments it is needed for formatting.
type ArgsInfo struct {
	Across    bool
	All       bool
	Blocks    bool
	Color     string
	Classify  bool
	Depth     uint8
	DirsFirst bool
	Filter    []*regexp.Regexp
	Group     bool
	Header    bool
	Ignore    []*regexp.Regexp
	Inode     bool
	Links     bool
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
