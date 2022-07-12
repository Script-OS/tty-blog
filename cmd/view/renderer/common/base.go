package common

import (
	"github.com/muesli/ansi"
	"tty-blog/cmd/view/renderer/style"
)

type BlockDecorator struct {
	Prefix     func(lineNo int, lineNum int) string
	Suffix     func(lineNo int, lineNum int) string
	LinePrefix string
	LineSuffix string
	InnerStyle style.Style
	IsThin     bool
	Used       bool
	width      *int
}

func (deco *BlockDecorator) Deco(line string, lineNo int, lineNum int) string {
	prefix := ""
	suffix := ""
	if deco.Prefix != nil {
		prefix = deco.Prefix(lineNo, lineNum)
	}
	if deco.Suffix != nil {
		suffix = deco.Suffix(lineNo, lineNum)
	}
	ret := prefix + line + suffix
	return ret
}

func (deco *BlockDecorator) Style() style.Style {
	if deco.InnerStyle == nil {
		return style.Style{}
	} else {
		return deco.InnerStyle
	}
}

func (deco *BlockDecorator) Push() string {
	return deco.LinePrefix
}

func (deco *BlockDecorator) Pop() string {
	return deco.LineSuffix
}

func (deco *BlockDecorator) Thin() bool {
	return deco.IsThin
}

func (deco *BlockDecorator) Width() int {
	if deco.width == nil {
		width := 0
		if deco.Prefix != nil {
			width += ansi.PrintableRuneWidth(deco.Prefix(0, 0))
		}
		if deco.Suffix != nil {
			width += ansi.PrintableRuneWidth(deco.Suffix(0, 0))
		}
		//width += ansi.PrintableRuneWidth(deco.OncePrefix)
		//width += ansi.PrintableRuneWidth(deco.OnceSuffix)
		deco.width = &width
	}
	return *deco.width
}
