package main

import (
	"regexp"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	colorNever  = "never"
	colorAlways = "always"
	colorAuto   = "auto"
)

type argsInfo struct {
	source    []string
	separator string
	across    bool
	all       bool
	color     string
	classify  bool
	depth     int
	filter    *regexp.Regexp
	ignore    *regexp.Regexp
	inode     bool
	long      bool
	oneLine   bool
	reverse   bool
	recursive bool
	tree      bool
}

func parseArgs() argsInfo {
	var args argsInfo
	kingpin.Arg("source", "the directory to list contents").Default(".").StringsVar(&args.source)
	kingpin.Flag("separator", "separator of the columns").Short('S').Default("  ").StringVar(&args.separator)
	kingpin.Flag("across", "writes the entries by lines instead of by columns").Short('x').BoolVar(&args.across)
	kingpin.Flag("all", "do not hide entries starting with .").Short('a').BoolVar(&args.all)
	kingpin.Flag("color", "control whether color is used to distinguish file types").Default(colorAuto).EnumVar(&args.color, colorNever, colorAlways, colorAuto)
	kingpin.Flag("classify", "append indicator to the entries").Short('F').BoolVar(&args.classify)
	kingpin.Flag("depth", "maximum depth of recursion").Short('D').IntVar(&args.depth)
	kingpin.Flag("filter", "only show entries that matches the pattern").Short('f').RegexpVar(&args.filter)
	kingpin.Flag("ignore", "do not show entries that matches the pattern").Short('I').RegexpVar(&args.ignore)
	kingpin.Flag("inode", "show entry inode").Short('i').BoolVar(&args.inode)
	kingpin.Flag("long", "show entries in the \"long view\"").Short('l').BoolVar(&args.long)
	kingpin.Flag("one-line", "show one entry per line").Short('1').BoolVar(&args.oneLine)
	kingpin.Flag("reverse", "reverse order of entries").Short('r').BoolVar(&args.reverse)
	kingpin.Flag("recursive", "list subdirectories recursively").Short('R').BoolVar(&args.recursive)
	kingpin.Flag("tree", "shows the entries in the tree view").Short('t').BoolVar(&args.tree)
	kingpin.Parse()
	return args
}
