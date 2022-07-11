package global

import "io/fs"

type CMD = func(args []string)

var User = "guest"

var Root fs.FS

var WorkDir string = "."
