package ls

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"io/fs"
	"os"
	"path/filepath"
	"tty-blog/global"
)

const Name = "ls"

func Run(args []string) {
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of ls:\n  ls [dir]")
		flagSet.PrintDefaults()
	}
	err := flagSet.Parse(args)
	if err == flag.ErrHelp {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}

	if flagSet.NArg() > 1 {
		fmt.Fprintln(os.Stderr, termenv.String("ls need at most one argument").Foreground(termenv.ANSIRed))
		return
	}
	dir := "."
	if flagSet.NArg() == 1 {
		dir = flagSet.Arg(0)
	}
	dir = global.CalcPath(filepath.Clean(dir))

	entries, err := fs.ReadDir(global.Root, dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}
	fileStyle := termenv.Style{}
	dirStyle := termenv.Style{}.Foreground(termenv.ANSIBlue)
	for _, entry := range entries {
		name := entry.Name()
		if name[0] == '.' {
			continue
		}
		if entry.IsDir() {
			fmt.Println(dirStyle.Styled(name))
		} else {
			fmt.Println(fileStyle.Styled(name))
		}
	}
}

var Completer = readline.PcItem(Name, global.NewPathCompleter())
