package ipe

import "regexp"

const (
	ArgColorNever  = "never"
	ArgColorAlways = "always"
	ArgColorAuto   = "auto"
)

type ArgsInfo struct {
	Across    bool
	All       bool
	Color     string
	Classify  bool
	Depth     uint8
	Filter    *regexp.Regexp
	Ignore    *regexp.Regexp
	Inode     bool
	Long      bool
	OneLine   bool
	Reverse   bool
	Recursive bool
	Separator string
	Sources   []string
	Tree      bool
	Width     int
}
