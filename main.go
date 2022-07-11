package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tty-blog/cmd/cd"
	"tty-blog/cmd/ls"
	"tty-blog/global"
)

func main() {
	global.Root = os.DirFS(".")

	RegisterCommand(ls.Name, ls.Run)
	RegisterCommand(cd.Name, cd.Run)

	reader, err := readline.NewEx(&readline.Config{
		AutoComplete: readline.NewPrefixCompleter(
			ls.Completer,
			cd.Completer,
		),
	})
	if err != nil {
		log.Panicln(err)
	}

	usernameStyle := termenv.Style{}.Bold().Foreground(termenv.ANSIGreen)
	pathStyle := termenv.Style{}.Bold().Foreground(termenv.ANSIBlue)
	for {
		fakePath := filepath.Clean("/" + global.WorkDir)
		reader.SetPrompt(usernameStyle.Styled(fmt.Sprintf("%s@%s", global.User, "blog")) + ":" + pathStyle.Styled(fakePath) + "> ")
		line, err := reader.Readline()
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			log.Panicln(err)
		}

		//act like real terminal
		tmpLine := strings.TrimSpace(line)
		if tmpLine == "\n" || tmpLine == "" {
			continue
		}

		parts := strings.Split(line, " ")
		args := []string{}
		for _, it := range parts {
			part := strings.TrimSpace(it)
			if part != "" {
				args = append(args, part)
			}
		}
		Dispatch(args)
	}
	return
}
