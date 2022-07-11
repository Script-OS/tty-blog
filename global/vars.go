package global

import (
	"io/fs"
	"path/filepath"
)

type CMD = func(args []string)

var User = "guest"

var Root fs.FS

var WorkDir string = "."

func CalcPath(path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Clean("/" + filepath.Join(WorkDir, path))
	}
	return filepath.Join(".", path)
}
