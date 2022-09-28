package ls

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/muesli/termenv"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"tty-blog/global"
)

const Name = "ls"

const timeLayout = "2006-01-02 15:04:05"

type EntrySorter []fs.DirEntry

func (a EntrySorter) Len() int      { return len(a) }
func (a EntrySorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a EntrySorter) Less(i, j int) bool {
	if a[i].IsDir() != a[j].IsDir() {
		return a[i].IsDir()
	}
	infoI, _ := a[i].Info()
	infoJ, _ := a[j].Info()
	timeI := infoI.ModTime()
	timeJ := infoJ.ModTime()
	return timeI.After(timeJ) || (timeI.Equal(timeJ) && a[i].Name() < a[j].Name())
}

func Run(args []string) {
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of ls:\n  ls [dir]")
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
		fmt.Fprintln(os.Stderr, termenv.String("ls need at most one argument").Foreground(termenv.ANSIRed))
		return
	}
	dir := "."
	if flagSet.NArg() == 1 {
		dir = flagSet.Arg(0)
	}
	dir = global.CalcPath(filepath.Clean(dir))

	entries, err := fs.ReadDir(global.Root, dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, termenv.String(err.Error()).Foreground(termenv.ANSIRed))
		return
	}
	fileStyle := termenv.Style{}
	dirStyle := termenv.Style{}.Foreground(termenv.ANSIBlue)
	sort.Sort(EntrySorter(entries))
	fmt.Println(" TIME               │ NAME ")
	fmt.Println("────────────────────┼──────")
	for _, entry := range entries {
		name := entry.Name()
		if name[0] == '.' {
			continue
		}
		info, _ := entry.Info()
		if entry.IsDir() {
			fmt.Println(info.ModTime().Format(timeLayout), "│", dirStyle.Styled(name))
		} else {
			fmt.Println(info.ModTime().Format(timeLayout), "│", fileStyle.Styled(name))
		}
	}
}

var Completer = readline.PcItem(Name, global.NewPathCompleter())
