package main

import (
	"os"

	"github.com/Nhanderu/ipe"
	"github.com/Nhanderu/trena"
	"github.com/alecthomas/kingpin"
)

func main() {
	a := parseArgs()
	s := ipe.NewFormatter(a)
	ss := s.String()
	os.Stdout.WriteString(ss)
}

func parseArgs() ipe.ArgsInfo {
	var args ipe.ArgsInfo
	kingpin.Arg("sources", "the directory to list contents").Default(".").StringsVar(&args.Sources)
	kingpin.Flag("separator", "separator of the columns").Short('S').Default("  ").StringVar(&args.Separator)
	kingpin.Flag("across", "writes the entries by lines instead of by columns").Short('x').BoolVar(&args.Across)
	kingpin.Flag("all", "do not hide entries starting with .").Short('a').BoolVar(&args.All)
	kingpin.Flag("color", "control whether color is used to distinguish file types").Default(ipe.ArgColorAuto).EnumVar(&args.Color, ipe.ArgColorNever, ipe.ArgColorAlways, ipe.ArgColorAuto)
	kingpin.Flag("classify", "append indicator to the entries").Short('F').BoolVar(&args.Classify)
	kingpin.Flag("depth", "maximum depth of recursion").Short('D').Uint8Var(&args.Depth)
	kingpin.Flag("filter", "only show entries that matches the pattern").Short('f').RegexpVar(&args.Filter)
	kingpin.Flag("ignore", "do not show entries that matches the pattern").Short('I').RegexpVar(&args.Ignore)
	kingpin.Flag("inode", "show entry inode").Short('i').BoolVar(&args.Inode)
	kingpin.Flag("long", "show entries in the \"long view\"").Short('l').BoolVar(&args.Long)
	kingpin.Flag("one-line", "show one entry per line").Short('1').BoolVar(&args.OneLine)
	kingpin.Flag("reverse", "reverse order of entries").Short('r').BoolVar(&args.Reverse)
	kingpin.Flag("recursive", "list subdirectories recursively").Short('R').BoolVar(&args.Recursive)
	kingpin.Flag("tree", "shows the entries in the tree view").Short('t').BoolVar(&args.Tree)
	kingpin.Parse()
	args.Width, _, _ = trena.Size()
	return args
}
