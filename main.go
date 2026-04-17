package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/termenv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tty-blog/cmd/cd"
	"tty-blog/cmd/edit"
	"tty-blog/cmd/ls"
	"tty-blog/cmd/su"
	"tty-blog/cmd/view"
	"tty-blog/global"
)

const banner = "" +
	"████████╗████████╗██╗   ██╗     ██████╗ ██╗      ██████╗  ██████╗ \n" +
	"╚══██╔══╝╚══██╔══╝╚██╗ ██╔╝     ██╔══██╗██║     ██╔═══██╗██╔════╝ \n" +
	"   ██║      ██║    ╚████╔╝█████╗██████╔╝██║     ██║   ██║██║  ███╗\n" +
	"   ██║      ██║     ╚██╔╝ ╚════╝██╔══██╗██║     ██║   ██║██║   ██║\n" +
	"   ██║      ██║      ██║        ██████╔╝███████╗╚██████╔╝╚██████╔╝\n" +
	"   ╚═╝      ╚═╝      ╚═╝        ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝ \n" +
	"                                                                  \n"

func main() {
	runewidth.EastAsianWidth = false
	runewidth.DefaultCondition.EastAsianWidth = false

	global.RealDir = *global.Config.RootDir
	global.Root = os.DirFS(global.RealDir)

	RegisterCommand(ls.Name, ls.Run)
	RegisterCommand(cd.Name, cd.Run)
	RegisterCommand(su.Name, su.Run)
	RegisterCommand(view.Name, view.Run)
	RegisterCommand(edit.Name, edit.Run)
	RegisterCommand("help", HelpCmd)
	RegisterCommand("?", func(args []string) {
		fmt.Println("Usable commands:", strings.Join([]string{ls.Name, cd.Name, su.Name, view.Name, edit.Name, "help"}, " "))
	})

	reader, err := readline.NewEx(&readline.Config{
		AutoComplete: readline.NewPrefixCompleter(
			ls.Completer,
			cd.Completer,
			su.Completer,
			view.Completer,
			edit.Completer,
			HelpCompleter(ls.Name, cd.Name, su.Name, view.Name, edit.Name, "help"),
			readline.PcItem("?"),
		),
	})
	if err != nil {
		log.Panicln(err)
	}

	bannerStyle := termenv.Style{}.Foreground(termenv.ANSICyan)
	fmt.Print(bannerStyle.Styled(banner))

	usernameStyle := termenv.Style{}.Bold().Foreground(termenv.ANSIGreen)
	pathStyle := termenv.Style{}.Bold().Foreground(termenv.ANSIBlue)
	for {
		termenv.SetWindowTitle("TTY-BLOG")
		fakePath := filepath.Clean("/" + global.WorkDir)
		reader.SetPrompt(usernameStyle.Styled(fmt.Sprintf("%s@%s", global.User, "blog")) + ":" + pathStyle.Styled(fakePath) + "> ")
		line, err := reader.Readline()
		if err == io.EOF || err == readline.ErrInterrupt {
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
}
