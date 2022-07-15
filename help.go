package main

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"os"
)

func HelpCmd(args []string) {
	flagSet := flag.NewFlagSet("help", flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of help:\n  help <cmd>")
		flagSet.PrintDefaults()
	}
	err := flagSet.Parse(args)
	if err == flag.ErrHelp {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}

	if flagSet.NArg() == 0 {
		flagSet.Usage()
		return
	}
	if flagSet.NArg() > 1 {
		fmt.Fprintln(os.Stderr, termenv.String("ls need at most one argument").Foreground(termenv.ANSIRed))
		return
	}

	Dispatch([]string{flagSet.Arg(0), "-help"})
}

func HelpCompleter(cmds ...string) *readline.PrefixCompleter {
	cmdTexts := []readline.PrefixCompleterInterface{}
	for _, cmd := range cmds {
		cmdTexts = append(cmdTexts, readline.PcItem(cmd))
	}
	return readline.PcItem("help", cmdTexts...)
}
