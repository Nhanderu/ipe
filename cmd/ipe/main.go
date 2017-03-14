package main

import (
	"os"
	"syscall"

	"github.com/Nhanderu/ipe/ipefmt"
	"github.com/Nhanderu/trena"
	. "github.com/alecthomas/kingpin"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(int(err.(syscall.Errno)))
	}
	os.Stdout.WriteString(ipefmt.NewFormatter(args).String())
}

func parseArgs() (ipefmt.ArgsInfo, error) {
	var args ipefmt.ArgsInfo

	Arg("sources", "the directory to list contents").
		Default(".").
		StringsVar(&args.Sources)

	Flag("across", "writes the entries by lines instead of by columns").
		Short('x').
		BoolVar(&args.Across)

	Flag("all", "do not hide entries starting with .").
		Short('a').
		BoolVar(&args.All)

	Flag("color", "control whether color is used to distinguish file types").
		Default(ipefmt.ArgColorAuto).
		EnumVar(&args.Color,
			ipefmt.ArgColorNever,
			ipefmt.ArgColorAlways,
			ipefmt.ArgColorAuto)

	Flag("classify", "append indicator to the entries").
		Short('F').
		BoolVar(&args.Classify)

	Flag("depth", "maximum depth of recursion").
		Short('D').
		Uint8Var(&args.Depth)

	Flag("dirs-first", "show directories first").
		BoolVar(&args.DirsFirst)

	Flag("filter", "only show entries that matches the pattern").
		Short('f').
		RegexpListVar(&args.Filter)

	Flag("group", "show group alongside user").
		Short('g').
		BoolVar(&args.Group)

	Flag("header", "show columns headers for long view").
		Short('h').
		BoolVar(&args.Header)

	Flag("ignore", "do not show entries that matches the pattern").
		Short('I').
		RegexpListVar(&args.Ignore)

	Flag("inode", "show entry inode").
		Short('i').
		BoolVar(&args.Inode)

	Flag("long", "show entries in the \"long view\"").
		Short('l').
		BoolVar(&args.Long)

	Flag("one-line", "show one entry per line").
		Short('1').
		BoolVar(&args.OneLine)

	Flag("reverse", "reverse order of entries").
		Short('r').
		BoolVar(&args.Reverse)

	Flag("recursive", "list subdirectories recursively").
		Short('R').
		BoolVar(&args.Recursive)

	Flag("sort", "field to sort by").
		Short('s').
		Default(ipefmt.ArgSortNone).
		EnumVar(&args.Sort,
			ipefmt.ArgSortNone,
			ipefmt.ArgSortInode,
			ipefmt.ArgSortMode,
			ipefmt.ArgSortSize,
			ipefmt.ArgSortAccessed,
			ipefmt.ArgSortModified,
			ipefmt.ArgSortCreated,
			ipefmt.ArgSortUser,
			ipefmt.ArgSortName)

	Flag("separator", "separator of the columns").
		Short('S').
		Default("  ").
		StringVar(&args.Separator)

	Flag("time", "define which timestamps to show").
		Short('T').
		Default(ipefmt.ArgTimeMod).
		EnumsVar(&args.Time,
			ipefmt.ArgTimeAcc,
			ipefmt.ArgTimeMod,
			ipefmt.ArgTimeCrt)

	Flag("tree", "shows the entries in the tree view").
		Short('t').
		BoolVar(&args.Tree)

	Parse()
	var err error
	args.Width, _, err = trena.Size()
	return args, err
}
