package view

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"tty-blog/global"
)

const Name = "view"

func Run(args []string) {
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of view:\n  view <markdown file>")
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

	filePath := filepath.Join(global.RealDir, file)
	s, err := os.Stat(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}
	if s.IsDir() {
		fmt.Fprintln(os.Stderr, termenv.String("view can not open dir").Foreground(termenv.ANSIRed))
		return
	}
	if path.Ext(file) != ".md" {
		fmt.Fprintln(os.Stderr, termenv.String("view can only open markdown file").Foreground(termenv.ANSIRed))
	}
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String("read file error").Foreground(termenv.ANSIRed))
	}
	marked, err := glamour.Render(string(f), "dark") //set marked as pager's content
	fmt.Fprintln(os.Stdout, termenv.String(marked).Foreground(termenv.ANSIRed))
	return
}
