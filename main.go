package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"io"
	"log"
	"os"
	"strings"
	"tty-blog/cmd/ls"
	"tty-blog/global"
)

func main() {
	global.Root = os.DirFS(".")

	RegisterCommand(ls.Name, ls.Run)

	reader, err := readline.NewEx(&readline.Config{
		AutoComplete: readline.NewPrefixCompleter(
			ls.Completer,
		),
	})
	if err != nil {
		log.Panicln(err)
	}

	usernameStyle := termenv.Style{}.Bold().Foreground(termenv.ANSIGreen)
	for {
		reader.SetPrompt(usernameStyle.Styled(fmt.Sprintf("%s@%s> ", global.User, "blog")))
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
		for i, it := range parts {
			parts[i] = strings.TrimSpace(it)
		}
		Dispatch(parts)
	}
	return
}
