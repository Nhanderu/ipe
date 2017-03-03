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
	Depth     int
	Filter    *regexp.Regexp
	Ignore    *regexp.Regexp
	Inode     bool
	Long      bool
	OneLine   bool
	Reverse   bool
	Recursive bool
	Separator string
	Source    []string
	Tree      bool
	Width     int
}
