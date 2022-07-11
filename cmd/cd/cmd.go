package cd

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

const Name = "cd"

func Run(args []string) {
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	err := flagSet.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}

	if flagSet.NArg() != 1 {
		fmt.Fprintln(os.Stderr, termenv.String("cd need one argument").Foreground(termenv.ANSIRed))
		return
	}
	dir := global.CalcPath(filepath.Clean(flagSet.Arg(0)))

	stat, err := fs.Stat(global.Root, dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String("no such dir").Foreground(termenv.ANSIRed))
		return
	}
	if !stat.IsDir() {
		fmt.Fprintln(os.Stderr, termenv.String("not a dir").Foreground(termenv.ANSIRed))
		return
	}
	global.WorkDir = filepath.Clean(dir)
}

var Completer = readline.PcItem(Name, global.NewPathCompleter())
