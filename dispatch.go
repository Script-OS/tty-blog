package main

import (
	"fmt"
	"github.com/muesli/termenv"
	"os"
	"tty-blog/global"
)

var aim = map[string]global.CMD{}

func RegisterCommand(name string, cmd global.CMD) {
	aim[name] = cmd
}

func Dispatch(args []string) {
	if cmd, ok := aim[args[0]]; ok {
		cmd(args[1:])
	} else {
		fmt.Fprintln(os.Stderr, termenv.String(fmt.Sprintf("unknown command `%s`", args[0])).Foreground(termenv.ANSIRed))
	}
}
