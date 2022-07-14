package view

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"golang.org/x/term"
	"io/fs"
	"os"
	"path/filepath"
	"tty-blog/cmd/view/renderer"
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

	w, _, _ := term.GetSize(int(os.Stdin.Fd()))
	r := renderer.New(w)
	rendered, _ := renderer.EasyRender(r, raw)
	//os.WriteFile("debug2.txt", rendered, 0777)
	fmt.Print("\x1b[?1000h")
	fmt.Print(string(rendered))
	//RenderInPage(string(rendered))
	fmt.Print("\x1b[?1000l")
}

var Completer = readline.PcItem(Name, global.NewPathCompleter())
