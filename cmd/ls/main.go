package ls

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"io/fs"
	"os"
	"tty-blog/global"
)

const Name = "ls"

func Run(args []string) {
	entries, err := fs.ReadDir(global.Root, global.WorkDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}
	fileStyle := termenv.Style{}
	dirStyle := termenv.Style{}.Foreground(termenv.ANSIBlue)
	for _, entry := range entries {
		name := entry.Name()
		if name[0] == '.' {
			continue
		}
		if entry.IsDir() {
			fmt.Println(dirStyle.Styled(name))
		} else {
			fmt.Println(fileStyle.Styled(name))
		}
	}
}

var Completer = readline.PcItem(Name, global.NewPathCompleter())
