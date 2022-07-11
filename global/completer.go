package global

import (
	"github.com/chzyer/readline"
	"io/fs"
	"strings"
)

func complete(base string) []string {
	ret := []string{}
	if base == "" {
		base = WorkDir
	}
	if base[len(base)-1] == '/' {
		base = base[:len(base)-1]
	}
	prefix := base
	if prefix == "." {
		prefix = ""
	} else {
		prefix += "/"
	}
	entries, _ := fs.ReadDir(Root, base)
	for _, entry := range entries {
		name := entry.Name()
		if name[0] == '.' {
			continue
		}
		if entry.IsDir() {
			ret = append(ret, prefix+name+"/")
		} else {
			ret = append(ret, prefix+name+" ")
		}
	}
	return ret
}

type PathCompleter struct {
	readline.PrefixCompleter
}

func (p *PathCompleter) GetDynamicNames(line []rune) [][]rune {
	parts := strings.Split(string(line), " ")
	base := parts[len(parts)-1]
	var names = [][]rune{}
	for _, name := range p.Callback(base) {
		names = append(names, []rune(name))
	}
	if len(names) == 0 {
		names = append(names, []rune{' '})
	}
	return names
}

func NewPathCompleter(pc ...readline.PrefixCompleterInterface) *PathCompleter {
	return &PathCompleter{
		PrefixCompleter: readline.PrefixCompleter{
			Callback: complete,
			Dynamic:  true,
			Children: pc,
		},
	}
}
