package edit

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"os"
	"os/exec"
	"path/filepath"
	"tty-blog/global"
)

const Name = "edit"

func Run(args []string) {
	if global.User != "editor" {
		fmt.Fprintln(os.Stderr, termenv.String("only editor can edit file").Foreground(termenv.ANSIRed))
		return
	}
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of edit:\n  edit <file>")
		flagSet.PrintDefaults()
	}
	err := flagSet.Parse(args)
	if err == flag.ErrHelp {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}

	if flagSet.NArg() != 1 {
		fmt.Fprintln(os.Stderr, termenv.String("edit need one argument").Foreground(termenv.ANSIRed))
		return
	}
	dir := global.CalcPath(filepath.Clean(flagSet.Arg(0)))

	editorArgs := append([]string{}, (*global.Config.Editor)[1:]...)
	editorArgs = append(editorArgs, filepath.Join(global.RealDir, dir))
	cmd := exec.Command((*global.Config.Editor)[0], editorArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

var Completer = readline.PcItem(Name, global.NewPathCompleter())
