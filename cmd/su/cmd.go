package su

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"os"
	"strings"
	"tty-blog/global"
)

const Name = "su"

func Run(args []string) {
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of su:\n  su [guest|editor]")
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
		fmt.Fprintln(os.Stderr, termenv.String("su need at most one argument").Foreground(termenv.ANSIRed))
		return
	}

	user := "editor"
	if flagSet.NArg() == 1 {
		user = strings.TrimSpace(flagSet.Arg(0))
	}
	if user == "editor" {
		if global.Config.EditorPassword != nil {
			reader, err := readline.New("")
			if err != nil {
				return
			}
			raw, err := reader.ReadPassword("input password:")
			if err != nil {
				return
			}
			if string(raw) != *global.Config.EditorPassword {
				fmt.Fprintln(os.Stderr, termenv.String("wrong password").Foreground(termenv.ANSIRed))
				return
			}
		} else {
			fmt.Fprintln(os.Stderr, termenv.String("WARN: editor doesn't has a password").Foreground(termenv.ANSIBrightYellow))
		}
		global.User = user
	} else if user == "guest" {
		global.User = user
	} else {
		fmt.Fprintln(os.Stderr, termenv.String("no such user").Foreground(termenv.ANSIRed))
	}
}

var Completer = readline.PcItem(Name, readline.PcItem("editor"), readline.PcItem("guest"))
