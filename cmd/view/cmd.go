package view

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"io/fs"
	"os"
	"path/filepath"
	"tty-blog/global"
)

const Name = "view"

func Run(args []string) {
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of view:\n  cd <file>")
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
		fmt.Fprintln(os.Stderr, termenv.String("view need one argument").Foreground(termenv.ANSIRed))
		return
	}
	file := global.CalcPath(filepath.Clean(flagSet.Arg(0)))

	stat, err := fs.Stat(global.Root, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String("no such file").Foreground(termenv.ANSIRed))
		return
	}
	if stat.IsDir() {
		fmt.Fprintln(os.Stderr, termenv.String("not a regular file").Foreground(termenv.ANSIRed))
		return
	}
	if filepath.Ext(file) != ".md" {
		fmt.Fprintln(os.Stderr, termenv.String("view can only open markdown file").Foreground(termenv.ANSIRed))
		return
	}

	raw, err := fs.ReadFile(global.Root, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}

	renderer, _ := glamour.NewTermRenderer(glamour.WithStylePath("dark"))
	rendered, _ := renderer.Render(string(raw))
	RenderInPage(rendered)
}

var Completer = readline.PcItem(Name, global.NewPathCompleter())
