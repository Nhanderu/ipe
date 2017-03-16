package main

import (
	"os"
	"syscall"

	"github.com/Nhanderu/ipe/ipefmt"
	"github.com/Nhanderu/trena"
	"github.com/alecthomas/kingpin"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(int(err.(syscall.Errno)))
	}
	_, err = ipefmt.NewFormatter(args).WriteTo(os.Stdout)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(int(err.(syscall.Errno)))
	}
}

func parseArgs() (ipefmt.ArgsInfo, error) {
	var args ipefmt.ArgsInfo

	kingpin.Arg("sources", "defines the directories to list contents").
		Default(".").
		StringsVar(&args.Sources)

	kingpin.Flag("across", "writes the entries by lines instead of by columns").
		Short('x').
		BoolVar(&args.Across)

	kingpin.Flag("all", "shows all entries").
		Short('a').
		BoolVar(&args.All)

	kingpin.Flag("blocks", "shows the number of file system blocks in long view").
		BoolVar(&args.Blocks)

	kingpin.Flag("color", "controls whether color is used").
		Default(ipefmt.ArgColorAuto).
		EnumVar(&args.Color,
			ipefmt.ArgColorNever,
			ipefmt.ArgColorAlways,
			ipefmt.ArgColorAuto)

	kingpin.Flag("classify", "appends indicator to the entries").
		Short('F').
		BoolVar(&args.Classify)

	kingpin.Flag("depth", "defines maximum depth of recursion").
		Short('D').
		Uint8Var(&args.Depth)

	kingpin.Flag("dirs-first", "shows directories first").
		BoolVar(&args.DirsFirst)

	kingpin.Flag("filter", "shows only the entries that matches the pattern").
		Short('f').
		RegexpListVar(&args.Filter)

	kingpin.Flag("group", "shows group alongside user").
		Short('g').
		BoolVar(&args.Group)

	kingpin.Flag("header", "shows columns headers in long view").
		Short('h').
		BoolVar(&args.Header)

	kingpin.Flag("ignore", "hides every entry that matches the pattern").
		Short('I').
		RegexpListVar(&args.Ignore)

	kingpin.Flag("inode", "shows entry inode in long view").
		Short('i').
		BoolVar(&args.Inode)

	kingpin.Flag("links", "shows the number of hard links in long view").
		BoolVar(&args.Links)

	kingpin.Flag("long", "display entries in \"long view\"").
		Short('l').
		BoolVar(&args.Long)

	kingpin.Flag("one-line", "shows one entry per line").
		Short('1').
		BoolVar(&args.OneLine)

	kingpin.Flag("reverse", "reverses order of entries").
		Short('r').
		BoolVar(&args.Reverse)

	kingpin.Flag("recursive", "lists subdirectories recursively").
		Short('R').
		BoolVar(&args.Recursive)

	kingpin.Flag("sort", "defines the field to sort by").
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

	kingpin.Flag("separator", "defines the separator of the columns").
		Short('S').
		Default("  ").
		StringVar(&args.Separator)

	kingpin.Flag("time", "defines which timestamps to show").
		Short('T').
		Default(ipefmt.ArgTimeMod).
		EnumsVar(&args.Time,
			ipefmt.ArgTimeAcc,
			ipefmt.ArgTimeMod,
			ipefmt.ArgTimeCrt)

	kingpin.Flag("tree", "display entries in \"tree view\"").
		Short('t').
		BoolVar(&args.Tree)

	kingpin.Parse()
	var err error
	args.Width, _, err = trena.Size()
	return args, err
}
